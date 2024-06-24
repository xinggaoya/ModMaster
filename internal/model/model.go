package model

/**
  @author: XingGao
  @date: 2024/6/24
**/

type GameInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Img  string `json:"img"`
}

type LocalGame struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
