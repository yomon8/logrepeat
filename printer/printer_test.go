package printer

func ExamplePrint() {
	Get().Run()
	Get().Spool <- "Test Message1"
	Get().Spool <- "Test Message2"
	// Output:
	// Test Message1
	// Test Message2
	// printer closed ..
	Close()
}
