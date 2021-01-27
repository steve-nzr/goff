package abstract

type Serializable interface {
	Serialize() (data []byte)
}

type Deserializable interface {
	Deserialize(data []byte)
}
