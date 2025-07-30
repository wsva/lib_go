package compress

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

func ZipCompressPath(sourcePath, destDirectory, destFilename string) error {
	sourcePathInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	err = mkdir(destDirectory)
	if err != nil {
		return err
	}

	zipFile, err := os.OpenFile(path.Join(destDirectory, destFilename),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipFileWriter := zip.NewWriter(zipFile)
	defer zipFileWriter.Close()

	if sourcePathInfo.IsDir() {
		baseDirectory := filepath.Base(sourcePath)
		_, err := zipFileWriter.Create(baseDirectory + "/")
		if err != nil {
			return err
		}
		return ZipWriterAddDirectory(zipFileWriter, sourcePath, baseDirectory)
	}
	return ZipWriterAddFile(zipFileWriter, sourcePath, "")
}

func ZipCompressPathList(sourcePathList []string, destDirectory, destFilename string) error {
	sourcePathInfoMap := make(map[string]os.FileInfo)
	for _, v := range sourcePathList {
		info, err := os.Stat(v)
		if err != nil {
			return err
		}
		sourcePathInfoMap[v] = info
	}

	err := mkdir(destDirectory)
	if err != nil {
		return err
	}

	zipFile, err := os.OpenFile(path.Join(destDirectory, destFilename),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipFileWriter := zip.NewWriter(zipFile)
	defer zipFileWriter.Close()

	for k, v := range sourcePathInfoMap {
		if v.IsDir() {
			err = ZipWriterAddDirectory(zipFileWriter, k, "")
			if err != nil {
				return err
			}
		} else {
			err = ZipWriterAddFile(zipFileWriter, k, v.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ZipWriterAddDirectory(zw *zip.Writer, sourceDirectory, directoryInZipfile string) error {
	fileInfoList, err := os.ReadDir(sourceDirectory)
	if err != nil {
		return err
	}
	for _, v := range fileInfoList {
		if v.IsDir() {
			_, err := zw.Create(directoryInZipfile + "/" + v.Name() + "/")
			if err != nil {
				return err
			}
			ZipWriterAddDirectory(zw,
				path.Join(sourceDirectory, v.Name()),
				directoryInZipfile+"/"+v.Name())
		} else {
			ZipWriterAddFile(zw,
				path.Join(sourceDirectory, v.Name()),
				directoryInZipfile+"/"+v.Name())
		}
	}
	return nil
}

func ZipWriterAddFile(zw *zip.Writer, sourceFilename, fullpathFilenameInZipfile string) error {
	fileContent, err := os.ReadFile(sourceFilename)
	if err != nil {
		return err
	}
	zipFile, err := zw.Create(fullpathFilenameInZipfile)
	if err != nil {
		return err
	}
	_, err = zipFile.Write(fileContent)
	if err != nil {
		return err
	}
	return nil
}

func ZipDecompressFile(zipFilename, destDirectory string) error {
	err := mkdir(destDirectory)
	if err != nil {
		return err
	}

	reader, err := zip.OpenReader(zipFilename)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		fileInfo := file.FileHeader.FileInfo()
		if fileInfo.IsDir() {
			fullpathDirectory := path.Join(destDirectory, file.Name)
			err = mkdir(fullpathDirectory)
			if err != nil {
				return err
			}
		} else {
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			fullpathFilename := path.Join(destDirectory, file.Name)
			w, err := os.Create(fullpathFilename)
			if err != nil {
				return err
			}
			defer w.Close()

			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func mkdir(directory string) error {
	if f, err := os.Stat(directory); err == nil {
		if f.IsDir() {
			return nil
		} else {
			return errors.New(directory + " exists and is not a directory")
		}
	}
	return os.MkdirAll(directory, 0777)
}
