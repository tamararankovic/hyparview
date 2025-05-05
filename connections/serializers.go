package connections

import "github.com/tamararankovic/hyparview/data"

func serialize(data.Message) ([]byte, error) {
	return nil, nil
}

func deserialize([]byte) (data.Message, error) {
	return data.Message{}, nil
}
