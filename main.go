package main

import (
	"archive/zip"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {

	gocron.Every(1).Hour().Do(task)

	<-gocron.Start()

}

func task() {

	mysqlBinDir := "C:/xampp71/mysql/bin"
	database := "adms_db"
	outputDir := "C:/goprojects"

	waktu := time.Now()

	baseName := fmt.Sprintf("%s-%d-%02d-%02d-%02d-%02d", database, waktu.Year(),
		waktu.Month(), waktu.Day(), waktu.Hour(), waktu.Minute())

	sqlOutputFile := fmt.Sprintf("%s/%s.sql", outputDir, baseName)

	zipOutputFile := fmt.Sprintf("%s/%s.zip", outputDir, baseName)

	mysqldump(mysqlBinDir, database, sqlOutputFile)

	zipFile(sqlOutputFile, zipOutputFile)

}

func zipFile(sourceFile string, targetFile string) {

	newZipFile, err := os.Create(targetFile)
	if err != nil {
		log.Fatal(err)
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	zipfile, err := os.Open(sourceFile)
	defer zipfile.Close()

	info, err := zipfile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	header, err := zip.FileInfoHeader(info)

	header.Name = filepath.Base(sourceFile)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)

	_, err = io.Copy(writer, zipfile)

}

func mysqldump(mysqlBinDir string, database string, sqlOutputFile string) string {
	args := []string{"-u", "root", database, "--ignore-table=adms_db.devlog"}
	cmd := exec.Command("mysqldump", args...)
	cmd.Dir = mysqlBinDir

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(stdout)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(sqlOutputFile, bytes, 0644)

	if err != nil {
		panic(err)
	}

	return sqlOutputFile

}
