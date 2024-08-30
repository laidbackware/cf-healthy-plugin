package sig_check

func sendIntNonBlock(targetChannel chan int, value int) {
	select {
	case targetChannel <- value:
	default:
	}
}

func getIntNonBlock(targetChannel chan int) int {
	select {
	case value := <- targetChannel:
		return value
	default:
		return 0
	}
}

func sendErrNonBlock(targetChannel chan error, value error) {
	select {
	case targetChannel <- value:
	default:
	}
}

func getErrNonBlock(targetChannel chan error) error {
	select {
	case value := <- targetChannel:
		return value
	default:
		return nil
	}
}