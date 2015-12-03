package haigo

import "io/ioutil"

func parseMongoFile(file string) (*HaigoFile, error) {

	b, err := ioutil.ReadFile(file)

	var hf HaigoFile

	err = hf.unmarshalYAML(b)
	if err != nil {
		return nil, err
	}

	return &hf, nil
}
