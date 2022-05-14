package commonutils

// SecurePanic only panic when err not nil
func SecurePanic(err error) {
	if err != nil {
		panic(err)
	}
}
