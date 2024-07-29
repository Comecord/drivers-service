package services

import "drivers-service/api/dto"

type Client struct {
	ID    string
	Name  string
	Value string
}

func (cl *Client) CreateClient(data dto.ClientCreate) {

}
