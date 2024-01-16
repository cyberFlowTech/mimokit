package lan

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
	"sync"
)

//go:embed locales
var LocaleFS embed.FS

var translateObj *Translate
var once sync.Once

type Translate struct {
	bundle    *i18n.Bundle
	localizer map[string]*i18n.Localizer
	lanList   []string
}

func (t *Translate) loadLan() error {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	files, err := LocaleFS.ReadDir("locales")
	if err != nil {
		logx.Errorf("ReadDir err : %s", err.Error())
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no locale files found")
	}

	for _, file := range files {
		fName := file.Name()
		path := fmt.Sprintf("locales/%s", fName)
		_, err := bundle.LoadMessageFileFS(LocaleFS, path)
		if err != nil {
			logx.Errorf("loadLan err : %s,path : %s", err.Error(), path)
			return err
		}
		t.lanList = append(t.lanList, fName[:len(fName)-5])
	}
	t.bundle = bundle
	return nil
}

func (t *Translate) newTranslator() {
	for _, v := range t.lanList {
		t.localizer[v] = i18n.NewLocalizer(t.bundle, t.getLanTag(v).String())
	}
}

func (t *Translate) getLanTag(lan string) language.Tag {
	lanTag, ok := SupportLan[lan]
	if !ok {
		return language.English
	}
	return lanTag
}

func InitTranslate() {
	once.Do(func() {
		translateObj = &Translate{}
		translateObj.localizer = make(map[string]*i18n.Localizer)
		translateObj.lanList = make([]string, 0)
		err := translateObj.loadLan()
		if err == nil {
			translateObj.newTranslator()
		}
	})
}

func GetTranslate() *Translate {
	if translateObj == nil {
		InitTranslate()
	}
	return translateObj
}

func Trans(lan string, key string, args ...interface{}) string {
	if p := GetTranslate(); p != nil {
		return p.Trans(lan, key, args)
	}
	//兜底文案
	return "The current network is congested, please wait [-0]"
}

func (t *Translate) getLocalizer(lan string) *i18n.Localizer {
	if t.localizer == nil || len(t.localizer) == 0 {
		return nil
	}

	localizerObj, ok := t.localizer[lan]
	if !ok {
		localizerObj = t.localizer["en"]
	}

	return localizerObj
}

func (t *Translate) Trans(lan string, key string, args ...interface{}) string {
	localizerObj := t.getLocalizer(lan)
	if localizerObj == nil {
		return "The current network is congested, please wait [-1]"
	}

	message, err := localizerObj.LocalizeMessage(&i18n.Message{ID: key})
	if err != nil {
		logx.Errorf("Message Trans err : %s, lan : %s, key : %s", err.Error(), lan, key)
		return "The current network is congested, please wait"
	}

	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}
	return message
}
