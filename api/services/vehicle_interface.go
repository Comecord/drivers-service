package services

import "drivers-service/api/dto"

type VehicleServiceApi interface {
	// Create creates a new vehicle post using the provided request.
	//
	// Parameters:
	// - post: The request containing the data for the new post.
	//
	// Returns:
	// - *dto.DBVehicle: The newly created post.
	// - error: An error if the creation operation fails.
	Create(*dto.CreateVehicleRequest) (*dto.DBVehicle, error)
	// Update updates a vehicle post with the given ID using the provided data.
	//
	// Parameters:
	// - id: The ID of the vehicle post to update.
	// - data: The updated data for the vehicle post.
	//
	// Returns:
	// - *dto.DBVehicle: The updated vehicle post.
	// - error: An error if the update operation fails.
	Update(string, *dto.UpdateVehicleRequest) (*dto.DBVehicle, error)
	// FindById retrieves a DBVehicle object based on the provided ID.
	//
	// Parameters:
	// - id: The ID of the post to retrieve.
	// Return type:
	// - *dto.DBVehicle: A pointer to the retrieved post.
	// - error: An error if the retrieval operation fails.
	FindById(string) (*dto.DBVehicle, error)
	// Find retrieves a list of DBVehicle objects based on the provided page and limit.
	//
	// Parameters:
	// - page: The page number of the results to retrieve. Defaults to 1 if not provided.
	// - limit: The maximum number of results to retrieve per page. Defaults to 10 if not provided.
	//
	// Returns:
	// - []*dto.DBVehicle: A slice of pointers to DBVehicle objects representing the retrieved posts.
	// - error: An error if the retrieval operation fails.
	Find(page int, limit int) ([]*dto.DBVehicle, error)
	// Delete deletes a post with the given ID.
	//
	// id: The ID of the post to be deleted.
	// Returns an error if the deletion fails.
	Delete(string) error
}
