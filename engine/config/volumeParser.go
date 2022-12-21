package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/swamp-labs/swamp/engine/volume"
	"gopkg.in/yaml.v3"
	"time"
)

//type volumeBlockTemplate []yaml.Node

type waitTemplate struct {
	Wait string `yaml:"wait" validate:"required"`
}

type constantTemplate struct {
	ConstantUserPerSec int    `yaml:"constantUserPerSec" validate:"required"`
	During             string `yaml:"during" validate:"required"`
}

type increaseTemplate struct {
	IncreaseUserPerSec variationTemplate `yaml:"increaseUserPerSec" validate:"required"`
	During             string            `yaml:"during" validate:"required"`
}

type decreaseTemplate struct {
	DecreaseUserPerSec variationTemplate `yaml:"decreaseUserPerSec" validate:"required"`
	During             string            `yaml:"during" validate:"required"`
}

type variationTemplate struct {
	from int `yaml:"from" validate:"required"`
	to   int `yaml:"to" validate:"required"`
}

type instantTemplate struct {
	instant int `yaml:"instant" validate:"required"`
}

func castTemplateNodeToVolume[T volumeTemplate](node yaml.Node, dst T) (*T, error) {
	err := node.Decode(&dst)
	if err != nil {
		return nil, fmt.Errorf("fail to cast node : %v", err)
	}
	validate := validator.New()
	err = validate.Struct(dst)
	if err != nil {
		return nil, fmt.Errorf("can not validate structure : %v", err)
	}
	return &dst, nil
}

func yamlNodeToVolume(node yaml.Node) (volume.Volume, error) {
	var dst volumeTemplate
	var err error

	dst, err = castTemplateNodeToVolume[waitTemplate](node, waitTemplate{})
	if err == nil {
		return dst.toVolume(), nil
	}
	dst, err = castTemplateNodeToVolume[constantTemplate](node, constantTemplate{})
	if err == nil {
		return dst.toVolume(), nil
	}
	dst, err = castTemplateNodeToVolume[increaseTemplate](node, increaseTemplate{})
	if err == nil {
		return dst.toVolume(), nil
	}
	dst, err = castTemplateNodeToVolume[decreaseTemplate](node, decreaseTemplate{})
	if err == nil {
		return dst.toVolume(), nil
	}
	return nil, fmt.Errorf("invalid volume type")
}

type volumeTemplate interface {
	toVolume() volume.Volume
}

func (tpl waitTemplate) toVolume() volume.Volume {
	duration, _ := time.ParseDuration(tpl.Wait)
	return volume.NewWaitVolume(duration)
}

func (tpl constantTemplate) toVolume() volume.Volume {
	duration, _ := time.ParseDuration(tpl.During)
	return volume.NewConstantVolume(tpl.ConstantUserPerSec, duration)
}

func (tpl increaseTemplate) toVolume() volume.Volume {
	duration, _ := time.ParseDuration(tpl.During)
	return volume.NewIncreaseVolume(tpl.IncreaseUserPerSec.from, tpl.IncreaseUserPerSec.to, duration)
}

func (tpl decreaseTemplate) toVolume() volume.Volume {
	duration, _ := time.ParseDuration(tpl.During)
	return volume.NewDecreaseVolume(tpl.DecreaseUserPerSec.from, tpl.DecreaseUserPerSec.to, duration)
}
