package plant

// import "fmt"

// var (
// 	globalPlantRegistry *PlantSpecificationRegistry
// )

// func init() {
// 	globalPlantRegistry = NewPlantSpecificationRegistry()
// }

// type PlantSpecificationFactory func(map[string]interface{}) (PlantSpecification, error)

// type PlantSpecificationRegistry struct {
// 	registry map[string]PlantSpecificationFactory
// }

// func NewPlantSpecificationRegistry() *PlantSpecificationRegistry {
// 	return &PlantSpecificationRegistry{
// 		registry: make(map[string]PlantSpecificationFactory),
// 	}
// }

// func (r *PlantSpecificationRegistry) Register(category string, factory PlantSpecificationFactory) {
// 	r.registry[category] = factory
// }

// func (r *PlantSpecificationRegistry) Get(category string, params map[string]interface{}) (PlantSpecification, error) {
// 	factory, ok := r.registry[category]
// 	if !ok {
// 		return nil, fmt.Errorf("unknown plant category: %s", category)
// 	}
// 	return factory(params)
// }

// func (r *PlantSpecificationRegistry) Validate(category string, params map[string]interface{}) error {
// 	factory, ok := r.registry[category]
// 	if !ok {
// 		return fmt.Errorf("unknown plant category: %s", category)
// 	}
// 	spec, err := factory(params)
// 	if err != nil {
// 		return err
// 	}
// 	return spec.Validate()
// }
