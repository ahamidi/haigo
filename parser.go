package haigo

import "io/ioutil"

func parseMongoFile(file string) (*File, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var hf File

	err = hf.unmarshalYAML(b)
	if err != nil {
		return nil, err
	}

	return &hf, nil
}
