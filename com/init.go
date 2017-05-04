package com

func init() {
	if err := CoInitialize(nil); err != nil {
		panic(err)
	}
}
