module reliable

go 1.17

replace sender => ./sender

replace receiver => ./receiver

require (
	receiver v0.0.0-00010101000000-000000000000 // indirect
	sender v0.0.0-00010101000000-000000000000 // indirect
)
