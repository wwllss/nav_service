package dao

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func init() {
	db := GetDB()
	_ = db.AutoMigrate(&Version{})
}

type Version struct {
	gorm.Model
	Os            string
	App           string
	Version       string
	LastVersionID uint
	Owner         string
}

func FindNewestVersion(os string, app string) Version {
	db := GetDB()
	nv := Version{}
	db.Last(&nv, "os=? and app=?", os, app)
	return nv
}

func FindVersion(os string, app string, version string) Version {
	db := GetDB()
	lv := Version{}
	version = fixVersion(version)
	err := db.Last(&lv, "os=? and app=? and version=?", os, app, version).Error
	if err == nil {
		return lv
	}
	err = db.Last(&lv, "os=? and app=?", os, app).Error
	nv := Version{
		Os:            os,
		App:           app,
		Version:       version,
		LastVersionID: lv.ID,
	}
	if nv.Before(lv.Version) {
		panic(errors.New("版本号错误，最新版本号为" + lv.Version))
	}
	db.Create(&nv)
	return nv
}

func fixVersion(version string) string {
	if !strings.HasSuffix(version, ".99") {
		return version
	}
	arr := strings.Split(version, ".")
	if len(arr) != 3 {
		return version
	}
	mid, _ := strconv.Atoi(arr[1])
	arr[1] = strconv.Itoa(mid + 1)
	arr[2] = "0"
	return strings.Join(arr, ".")
}

func FindVersionById(id uint) Version {
	db := GetDB()
	v := Version{}
	db.Last(&v, "id=?", id)
	return v
}

func (v Version) FindLastVersion() Version {
	return FindVersionById(v.LastVersionID)
}

func (v Version) After(ver string) bool {
	return CompareVersion(v.Version, ver) == 1
}

func (v Version) Before(ver string) bool {
	return CompareVersion(v.Version, ver) == -1
}

func CompareVersion(version1, version2 string) int {
	n, m := len(version1), len(version2)
	i, j := 0, 0
	for i < n || j < m {
		x := 0
		for ; i < n && version1[i] != '.'; i++ {
			x = x*10 + int(version1[i]-'0')
		}
		i++
		y := 0
		for ; j < m && version2[j] != '.'; j++ {
			y = y*10 + int(version2[j]-'0')
		}
		j++
		if x > y {
			return 1
		}
		if x < y {
			return -1
		}
	}
	return 0
}
