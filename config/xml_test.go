package config

import (
	"fmt"
	"testing"
)

const (
	xml1 = `
		<server>
			<vars>
				<item name="ROOT">/var/www/</item>
			</vars>

			<net>
				<address>localhost</address>
				<port>9000</port>
				<protocol>tcp</protocol>
			</net>

			<user>
				<page path="/">
					<title>My page</title>
					<serve>%ROOT%/index.html</serve>
				</page>
			</user>
		</server>
	`
)

func TestXml1(t *testing.T) {
	conf, err := GetUserConfFromXml([]byte(xml1))

	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(string(conf))
	}
}

func TestXml2(t *testing.T) {
	vars, err := GetVarsFromXml([]byte(xml1))

	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(vars.ToMap())
	}
}
func TestXml3(t *testing.T) {
	conf, err := GetServerConfFromXml([]byte(xml1))

	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(conf)
	}
}
