package main

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type GameInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// GetGameList returns a greeting for the given name
func (a *App) GetGameList(name string) []GameInfo {
	url := "https://flingtrainer.com/?s=" + name
	c := colly.NewCollector()

	var gameList []GameInfo
	// 搜索出链接
	c.OnHTML(".post-content h2 a", func(e *colly.HTMLElement) {
		fmt.Printf("%s %s\n", e.Text, e.Attr("href"))
		title := e.Text
		// 正则title的空格字符转为下划线
		title = strings.Replace(title, " ", "_", -1)
		gameList = append(gameList, GameInfo{
			Name: title,
			Url:  e.Attr("href"),
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)

	return gameList
}

// GetGameInfo 查询游戏详情
func (a *App) GetGameInfo(url string) GameInfo {
	c := colly.NewCollector()
	var info GameInfo
	c.OnHTML(".zip .attachment-title a", func(e *colly.HTMLElement) {
		if info.Url == "" {
			info.Url = e.Attr("href")
			info.Name = e.Text
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
	DownloadGame(info)
	return info
}

type localGame struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// GetGame 获取游戏列表
func (a *App) GetGame() []localGame {
	var gameList []localGame
	if err := filepath.Walk("./execute", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		gameList = append(gameList, localGame{
			Name: info.Name(),
			Path: path,
		})
		return nil
	}); err != nil {
		log.Printf("遍历文件夹失败: %v", err)
	}
	return gameList
}

// RunGame 运行游戏
func (a *App) RunGame(path string) {
	// 获取程序运行路径
	execPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取程序运行路径失败: %v", err)
	}
	log.Printf("程序运行路径: %s", path)
	path = execPath + "\\" + path
	// 运行游戏 路径加引号
	cmd := exec.Command(path)
	err = cmd.Run()
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
func DownloadGame(info GameInfo) {
	res, err := http.Get(info.Url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// 判断创建文件夹
	if err := os.MkdirAll("./download", os.ModePerm); err != nil {
		log.Printf("创建文件夹失败: %v", err)
	}
	// 保存文件
	file, err := os.Create("./download/" + info.Name + ".zip")
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}(file)
	if err != nil {
		log.Printf("创建文件失败: %v", err)
	}
	io.Copy(file, res.Body)

	// 压缩包地址
	zipPath := "./download/" + info.Name + ".zip"
	// 解压
	err = Unzip(zipPath, "./execute/")
	if err != nil {
		log.Printf("解压失败: %v", err)
	}

	defer func(name string) {
		err = os.RemoveAll("download")
		if err != nil {
			log.Printf("删除文件失败: %v", err)
		}
	}(zipPath)
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
