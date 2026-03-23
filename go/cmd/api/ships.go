package main

type ShipService interface {
	ListShips() []string
}

type staticShipService struct {
	ships []string
}

func newShipService() ShipService {
	return staticShipService{
		ships: []string{
			"Millennium Falcon",
			"X-wing",
			"TIE Fighter",
			"Star Destroyer",
			"Naboo Royal Starship",
		},
	}
}

func (s staticShipService) ListShips() []string {
	return append([]string(nil), s.ships...)
}
