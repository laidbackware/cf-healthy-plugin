package sig_check

func getIntNonBlock(targetChannel chan int) int {
	select {
	case value := <-targetChannel:
		return value
	default:
		return 0
	}
}

func getErrNonBlock(targetChannel chan error) error {
	select {
	case value := <-targetChannel:
		return value
	default:
		return nil
	}
}
