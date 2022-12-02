package depinit

import (
	"container/list"
	"fmt"
)

type ErrCyclicDependency struct {
	modA string
	modB string
}

func (e *ErrCyclicDependency) Error() string {
	return fmt.Sprintf("cyclic dependency detected between (%s and %s)", e.modA, e.modB)
}

type DepMngr struct {
	allowCyclicDep bool

	modFunc map[string]func() error
	modDeps map[string][]string
}

func NewDepMngr(allowCyclicDep bool) *DepMngr {
	return &DepMngr{
		allowCyclicDep: allowCyclicDep,

		modFunc: make(map[string]func() error),
		modDeps: make(map[string][]string),
	}
}

// AddModule adds a new module to be initialized later.
// name: the unique identifier for this module.
// f: the function to initialize this module.
// deps: the unique identifiers of the modules this module depends on.
func (m *DepMngr) AddModule(name string, f func() error, deps ...string) {
	if _, ok := m.modFunc[name]; ok {
		// Module already exists
		return
	}

	m.modFunc[name] = f
	m.modDeps[name] = deps
}

// Init initializes all the registered modules in a dependency-aware order.
// This means that all the dependencies of a module will always be initialized before it.
// If a cyclic dependency exists and allowCyclicDep is false, Init will return ErrCyclicDependency.
func (m *DepMngr) Init() error {
	// Evaluate inverse dependencies and dependency count for topological sort
	modInvDeps := make(map[string][]string)
	modDepCount := make(map[string]int)
	for mod, deps := range m.modDeps {
		for _, dep := range deps {
			modInvDeps[dep] = append(modInvDeps[dep], mod)
		}
		modDepCount[mod] = len(deps)
	}

	// Initialize queue with root modules (modules with no dependency)
	queue := list.New()
	for mod := range m.modFunc {
		if modDepCount[mod] == 0 {
			queue.PushBack(mod)
		}
	}

	// Run topological sort
	for queue.Len() != 0 {
		cur := queue.Remove(queue.Front()).(string)

		if f, ok := m.modFunc[cur]; ok {
			// Init this mod
			if err := f(); err != nil {
				return err
			}
		}

		// Reduce dep count of subsequent mods
		for _, invDep := range modInvDeps[cur] {
			modDepCount[invDep]--
			if modDepCount[invDep] == 0 {
				// No more pending dependency for invDep, add to queue
				queue.PushBack(invDep)
			}
		}
	}

	// TODO: parallelize independent inits

	// TODO: check for cyclic dependencies

	return nil
}
