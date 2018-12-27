package fetcher

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type WikiDownloader struct {
	DownloaderName         string
	DownloadOutputTemplate string

	TmpDownloadDir   string
	DownloadNo       int
	DownResourcePath string
}

func NewDefaultWikiDownloader() *WikiDownloader {
	return &WikiDownloader{
		DownloaderName:         "wget",
		DownloadOutputTemplate: "/data/crawler/wiki/",
	}
}

func NewWikiDownloader(downloaderName, outputTemplate, tmpDownloadDir string, downloadNo int) *WikiDownloader {
	return &WikiDownloader{
		DownloaderName:         downloaderName,
		DownloadOutputTemplate: outputTemplate,
		TmpDownloadDir:         tmpDownloadDir,
		DownloadNo:             downloadNo,
	}
}

func (dl *WikiDownloader) SetTmpDownloadDir(tmpDownloadDir string) {
	dl.TmpDownloadDir = tmpDownloadDir
}

func (dl *WikiDownloader) SetDownloadNo(downloadNo int) {
	dl.DownloadNo = downloadNo
}

func (dl *WikiDownloader) Download(url, name string) ([]byte, error) {
	if dl.TmpDownloadDir == "" {
		return nil, fmt.Errorf("TmpDownloadDir field must not be null")
	}

	//today := time.Now().Format("2006-01-02")
	dl.DownResourcePath = path.Join(dl.DownloadOutputTemplate, dl.TmpDownloadDir)

	cmd := exec.Command("sh", "-c", "mkdir -p "+dl.DownResourcePath)
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("create directory: %s error: %s", dl.DownResourcePath, err)
	}

	// wget -q --show-progress --page-requisites --html-extension --convert-links --random-wait -e robots=off -nd --span-hosts $URL || true
	wgetCmd := fmt.Sprintf("%s -q --show-progress --page-requisites --html-extension --convert-links --random-wait -e robots=off -nd --span-hosts '%s' || true", dl.DownloaderName, url)

	fmt.Println("Download: ", wgetCmd)

	cmd = exec.Command("sh", "-c", wgetCmd)
	cmd.Dir = dl.DownResourcePath

	_, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("wait %s error: %s", dl.DownloaderName, err)
	}

	// 重命名html
	pageNameCmd := fmt.Sprintf("ls -S | grep -a .html | head -n1")
	cmd = exec.Command("sh", "-c", pageNameCmd)
	cmd.Dir = dl.DownResourcePath
	pageNameBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Can't not get page name of wiki webpage")
	}

	pageName := strings.Replace(string(pageNameBytes), "\n", "", -1)

	renameCmd := fmt.Sprintf("mv %s index.html", pageName)
	cmd = exec.Command("sh", "-c", renameCmd)
	cmd.Dir = dl.DownResourcePath
	_, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Can't rename %s to index.html", pageName)
	}

	body, err := ioutil.ReadFile(path.Join(dl.DownResourcePath, "index.html"))
	if err != nil {
		return nil, fmt.Errorf("Can't read %s contents", path.Join(dl.DownResourcePath, "index.html"))
	}

	return body, nil
}

func (dl *WikiDownloader) Size() (int, error) {
	commandStr := fmt.Sprintf("du -d 0 -k %s", dl.DownResourcePath)
	cmd := exec.Command("sh", "-c", commandStr)
	output, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("du command error: %s", err)
	}

	str := strings.Split(string(output), "\t")
	if len(str) == 0 {
		return -1, fmt.Errorf("can't get file size")
	}

	size, err := strconv.Atoi(str[0])
	if err != nil {
		return -1, fmt.Errorf("atoi error: %s, str: %s", err, str[0])
	}

	return size, nil
}

func (dl *WikiDownloader) GetDownResourcePath() string {
	return dl.DownResourcePath
}
