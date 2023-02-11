package jsonld

func Actor(reader Reader) Reader {
	return reader.Get("actor")
}

func Activity(reader Reader) Reader {
	return reader.Get("activity")
}

func Object(reader Reader) Reader {
	return reader.Get("object")
}
