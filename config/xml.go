package config

import (
	"encoding/xml"

	"github.com/s2ks/fcgiserver/util"
)

/* Defines structures to unmarshal from xml config file */

type VarXml struct {
	Items []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:",innerxml"`
	} `xml:"item"`
}

type NetXml struct {
	Address  string `xml:"address"`
	Port     string `xml:"port"`
	Protocol string `xml:"protocol"`
}

type XmlConf struct {
	XMLName xml.Name `xml:"server"`
	Net     NetXml   `xml:"net"`
}

type VarsXml struct {
	XMLName xml.Name `xml:"server"`
	Items   []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:",innerxml"`
	} `xml:"vars>item"`
}

type UserXml struct {
	XMLName xml.Name `xml:"user"`
	Raw     []byte   `xml:",innerxml"`
}

type UserConf struct {
	XMLName xml.Name `xml:"server"`
	User    UserXml
}

func GetVarsFromXml(raw []byte) (*VarsXml, error) {
	vars := new(VarsXml)

	err := xml.Unmarshal(raw, vars)

	if err != nil {
		return nil, err
	} else {
		return vars, nil
	}
}

func GetVarsFromXmlFile(path string) (*VarsXml, error) {
	buf, err := util.ReadFromFile(path)

	if err != nil {
		return nil, err
	}

	return GetVarsFromXml(buf)
}

func (vars *VarsXml) ToMap() map[string]string {
	varmap := make(map[string]string)

	for _, item := range vars.Items {
		varmap[item.Name] = item.Value
	}

	return varmap
}

func GetServerConfFromXml(raw []byte) (*XmlConf, error) {
	config := new(XmlConf)

	vars, err := GetVarsFromXml(raw)

	if err != nil {
		return nil, err
	}

	buf, err := util.ByteSubstituteMap(raw, vars.ToMap(), "%")

	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(buf, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetServerConfFromXmlFile(path string) (*XmlConf, error) {
	buf, err := util.ReadFromFile(path)

	if err != nil {
		return nil, err
	}

	return GetServerConfFromXml(buf)
}

/*
	TODO support for an <include/> tag
*/
func GetUserConfFromXml(raw []byte) ([]byte, error) {
	vars, err := GetVarsFromXml(raw)

	if err != nil {
		return nil, err
	}

	buf, err := util.ByteSubstituteMap(raw, vars.ToMap(), "%")

	if err != nil {
		return nil, err
	}

	userConf := new(UserConf)

	/* Get the raw inner xml of <user> */
	err = xml.Unmarshal(buf, userConf)

	if err != nil {
		return nil, err
	}

	/* Re-add <user> tags */
	dest, err := xml.Marshal(userConf.User)

	if err != nil {
		return nil, err
	}

	return dest, nil
}

func GetUserConfFromXmlFile(path string) ([]byte, error) {
	buf, err := util.ReadFromFile(path)

	if err != nil {
		return nil, err
	}

	return GetUserConfFromXml(buf)
}
