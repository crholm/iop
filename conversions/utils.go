package conversions

import (
	"encoding/json"
	"encoding/xml"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"io"
)

func decoderJSON(r io.Reader) decoder {
	return json.NewDecoder(r)
}
func encoderJSON(w io.Writer) encoder {
	return json.NewEncoder(w)
}

func encoderXML(w io.Writer) encoder {
	return xml.NewEncoder(w)
}
func decoderXML(r io.Reader) decoder {
	return xml.NewDecoder(r)
}

func decoderTOML(r io.Reader) decoder {
	return toml.NewDecoder(r)
}

func encoderTOML(w io.Writer) encoder {
	return toml.NewEncoder(w)
}

func decoderYAML(r io.Reader) decoder {
	return yaml.NewDecoder(r)
}

func encoderYAML(w io.Writer) encoder {
	return yaml.NewEncoder(w)
}
