package domain

type MinimaxAlg struct{}

func (alg MinimaxAlg) Score(game Game, depth int) int {
	if game.Win(int(game.GetPlayerChar())) {
		return 10 - depth
	} else if game.Win(int(game.GetOpponentChar())) {
		return depth - 10
	}
	return 0
}

func (alg MinimaxAlg) Minimax(game Game, depth int) (int, int) {
	if game.GameIsOver() {
		return alg.Score(game, depth), 0
	}

	depth++
	scores := make([]int, 0)
	moves := make([]int, 0)

	// Получаем доступные ходы
	availableMoves := game.GetAvailableMoves()

	// Для каждого возможного хода вычисляем оценку
	for _, move := range availableMoves {
		possibleGame := game.GetNewState(move)
		score, _ := alg.Minimax(possibleGame, depth)
		scores = append(scores, score)
		moves = append(moves, move)
	}

	// Выбираем минимальный или максимальный результат
	if game.GetActiveTurn() == game.GetPlayerChar() {
		// Максимизирующий игрок
		maxScoreIndex := 0
		maxScore := scores[0]
		for i := 1; i < len(scores); i++ {
			if scores[i] > maxScore {
				maxScore = scores[i]
				maxScoreIndex = i
			}
		}
		return maxScore, moves[maxScoreIndex]
	} else {
		// Минимизирующий игрок
		minScoreIndex := 0
		minScore := scores[0]
		for i := 1; i < len(scores); i++ {
			if scores[i] < minScore {
				minScore = scores[i]
				minScoreIndex = i
			}
		}
		return minScore, moves[minScoreIndex]
	}
}

func (alg MinimaxAlg) GetNextMove(game Game) (int, int) {
	score, choice := alg.Minimax(game, 0)
	return score, choice
}
