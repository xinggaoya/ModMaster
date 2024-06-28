package internal

import (
	"ModMaster/internal/consts"
	"ModMaster/internal/model"
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	go CheckUpdate()
}

// GetGameList returns a greeting for the given name
func (a *App) GetGameList(name string) []model.GameInfo {
	url := "https://flingtrainer.com/?s=" + name
	c := colly.NewCollector()

	var gameList []model.GameInfo
	c.OnHTML("article", func(e *colly.HTMLElement) {
		var info model.GameInfo
		e.ForEach(".post-content h2 a", func(i int, e *colly.HTMLElement) {
			info.Name = e.Text
			info.Url = e.Attr("href")
		})
		e.ForEach(".post-details .post-details-thumb a img", func(i int, e *colly.HTMLElement) {
			info.Img = e.Attr("src")
		})
		gameList = append(gameList, info)
	})
	c.Visit(url)

	return gameList
}

// GetGameListPage 获取游戏分页列表
func (a *App) GetGameListPage(page int) []model.GameInfo {
	url := "https://flingtrainer.com/page/" + fmt.Sprintf("%d", page)
	var gameList []model.GameInfo
	c := colly.NewCollector()
	c.OnHTML("article", func(e *colly.HTMLElement) {
		var info model.GameInfo
		e.ForEach(".post-content h2 a", func(i int, e *colly.HTMLElement) {
			info.Name = e.Text
			info.Url = e.Attr("href")
		})
		e.ForEach(".post-details .post-details-thumb a img", func(i int, e *colly.HTMLElement) {
			info.Img = e.Attr("src")
		})
		gameList = append(gameList, info)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
	return gameList
}

// GetGameInfo 查询游戏详情
func (a *App) GetGameInfo(url string, img string) model.LocalGame {
	c := colly.NewCollector()
	var info model.GameInfo
	c.OnHTML(".zip .attachment-title a", func(e *colly.HTMLElement) {
		if info.Url == "" {
			info.Url = e.Attr("href")
			info.Name = e.Text
		}
	})
	c.Visit(url)
	exePath := DownloadGame(info)
	return model.LocalGame{
		Img:  img,
		Name: info.Name,
		Path: exePath,
	}
}

// RunGame 运行游戏
func (a *App) RunGame(path string) {
	// 获取程序运行路径
	execPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取程序运行路径失败: %v", err)
	}
	path = filepath.Join(execPath, path)
	log.Printf("程序运行路径: %s", path)

	// 运行游戏 路径加引号
	cmd := exec.Command(path)
	err = cmd.Start()
	if err != nil {
		log.Printf("运行游戏失败: %v", err)
	}
}

// DeleteGame 刪除游戏
func (a *App) DeleteGame(path string) {
	// 删除文件
	if err := os.RemoveAll(path); err != nil {
		log.Printf("删除文件失败: %v", err)
	}
}

// DownloadGame 根据url下载游戏
func DownloadGame(info model.GameInfo) string {
	// 压缩包地址
	zipPath := "./download/" + info.Name + ".zip"

	// 判断文件夹存在
	err := os.MkdirAll("./download", os.ModePerm)
	if err != nil {
		log.Printf("创建文件夹失败: %v", err)
	}
	err = DownloadFile(info.Url, zipPath)
	if err != nil {
		log.Printf("下载失败: %v", err)
		return ""
	}

	// 解压路径
	thePathToDecompressTheDecompression := "./execute/" + info.Name
	err = Unzip(zipPath, thePathToDecompressTheDecompression)
	if err != nil {
		log.Printf("解压失败: %v", err)
	}

	//defer func(name string) {
	//	err = os.RemoveAll("./download")
	//	if err != nil {
	//		log.Printf("删除文件失败: %v", err)
	//	}
	//}(zipPath)

	var exePath string
	// 获取该路径下解压出来的exe
	filepath.Walk(thePathToDecompressTheDecompression, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("遍历文件夹失败: %v", err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".exe" {
			// 获取路径
			exePath = path
		}
		return nil
	})

	return exePath
}

// Unzip 解压ZIP文件到指定目录
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func(r *zip.ReadCloser) {
		err = r.Close()
		if err != nil {
			log.Printf("关闭zip文件失败: %v", err)
		}
	}(r)

	for _, f := range r.File {
		err = extractAndWriteFile(f, dest)
		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", f.Name, err)
		}
	}

	return nil
}

// extractAndWriteFile 从zip.Reader中提取单个文件并写入到磁盘
func extractAndWriteFile(file *zip.File, dest string) error {
	path := filepath.Join(dest, file.Name)

	// 检查此路径是否为目录
	if file.FileInfo().IsDir() {
		err := os.MkdirAll(path, file.Mode())
		if err != nil {
			return err
		}
		return nil
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer func(srcFile io.ReadCloser) {
		err = srcFile.Close()
		if err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}(srcFile)

	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, srcFile)
	return err
}

// CheckUpdate 检查更新
func CheckUpdate() {
	url := "https://api.github.com/repos/xinggaoya/ModMaster/releases/latest"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		// 获取最新版本号
		var data struct {
			TagName string `json:"tag_name"`
			Assets  []struct {
				BrowserDownloadUrl string `json:"browser_download_url"`
			} `json:"assets"`
		}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			log.Printf("解析json失败: %v", err)
			return
		}
		// 比较版本号
		if data.TagName != consts.AppVersion {
			log.Printf("有新版本: %s", data.TagName)
			// 下载新版本
			err = DownloadFile(data.Assets[0].BrowserDownloadUrl, "./update.exe")
			if err != nil {
				log.Printf("下载新版本失败: %v", err)
			}
			// 运行新版本
			cmd := exec.Command("./update.exe")
			err = cmd.Start()
			if err != nil {
				log.Printf("运行新版本失败: %v", err)
			}
		} else {
			log.Printf("当前已是最新版本")
		}
	}
}

// DownloadFile 下载文件并显示进度条
func DownloadFile(url string, filePath string) error {
	log.Printf("Downloading %s", url)
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// set cookie
	req.Header.Set("Cookie", "ikgxPfc-S=hGER74M6uxkL9; oYZQOsHLcUbd=EXbUmG.oJ; QcGqsjwyzIVOhUJm=dftYO%5DRxpH; dzXSbtHy=peDNBPm")
	// 模拟浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 获取文件大小
	fileSize := resp.ContentLength

	// 创建进度条
	log.Printf("Downloading %s (%d bytes)\n", filePath, fileSize)
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 进度记录
	var downloaded int64
	for {
		// 读取数据
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// 写入文件
		_, err = file.Write(buf[:n])
		if err != nil {
			return err
		}

		// 更新进度
		downloaded += int64(n)

		log.Printf("Downloaded %d bytes (%.2f%%)\n", downloaded, float64(downloaded)/float64(fileSize)*100)
	}
	return nil
}
