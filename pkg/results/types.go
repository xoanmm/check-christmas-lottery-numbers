package results

// Number struct which information about
// each lottery number to check
type Number struct {
	Number    int    `json:"number"`
	BetAmount int    `json:"bet_amount"`
	Origin    string `json:"origin"`
}

// PersonNumbers struct which contains information
// about the numbers to check for a specific Person
type PersonNumbers struct {
	Numbers	[]Number `json:"numbers"`
	Owner string `json:"owner"`
}

// PersonNumbersToCheck struct which contains information
// about the numbers to check for a specific Person
type PersonNumbersToCheck struct {
	PersonsNumbers	[]PersonNumbers `json:"numbers_to_check"`
}


// Result struct which contain information about
// each lottery number result obtained from API
type Result struct {
	Numero     int `json:"numero"`
	Premio     int `json:"premio"`
	Timestsamp int `json:"timestamp"`
	Status     int `json:"status"`
	Error      int `json:"error"`
}

// NewResult allows to create a Result type struct providing all the information for it
func NewResult(numero int, premio int, timestamp int, status int, error int) *Result {
	return &Result{numero,  premio, timestamp, status, error}
}

// LotteryDrawStatus which contain information about
// the actual status of the lottery draw
type LotteryDrawStatus struct {
	Status int `json:"status"`
	Error  int `json:"error"`
}

// NewLotteryDrawStatus allows to create a DrawStatus type struct providing all the information for it
func NewLotteryDrawStatus(status int, error int) *LotteryDrawStatus {
	return &LotteryDrawStatus{status, error}
}

// NumberResult struct which contain information
// about each lottery number result
type NumberResult struct {
	Number       int    `json:"number"`
	AmountEarned int    `json:"amount_earned"`
	Origin       string `json:"origin"`
}

// NumbersResult struct which contains
// an array of NumberResult
type NumbersResult struct {
	NumberResults []NumberResult `json:"numbers_result"`
}
