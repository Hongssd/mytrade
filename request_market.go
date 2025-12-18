package mytrade

type KlineParam struct {
	AccountType string
	Symbol      string
	Interval    string
	StartTime   int64
	EndTime     int64
	Limit       int
}

func (k *KlineParam) SetAccountType(accountType string) *KlineParam {
	k.AccountType = accountType
	return k
}

func (k *KlineParam) SetSymbol(symbol string) *KlineParam {
	k.Symbol = symbol
	return k
}

func (k *KlineParam) SetInterval(interval string) *KlineParam {
	k.Interval = interval
	return k
}

func (k *KlineParam) SetStartTime(startTime int64) *KlineParam {
	k.StartTime = startTime
	return k
}

func (k *KlineParam) SetEndTime(endTime int64) *KlineParam {
	k.EndTime = endTime
	return k
}

func (k *KlineParam) SetLimit(limit int) *KlineParam {
	k.Limit = limit
	return k
}

type BookParam struct {
	AccountType string
	Symbol      string
	Level       int

	// Sunx
	Step string
}

func (b *BookParam) SetAccountType(accountType string) *BookParam {
	b.AccountType = accountType
	return b
}

func (b *BookParam) SetSymbol(symbol string) *BookParam {
	b.Symbol = symbol
	return b
}

func (b *BookParam) SetLevel(level int) *BookParam {
	b.Level = level
	return b
}
