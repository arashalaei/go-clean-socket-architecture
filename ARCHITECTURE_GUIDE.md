# Clean Architecture: Entities vs Repository Models

## Recommended Approach: **Separate with Mappers** ✅

### Why Separate?

1. **Independence**: Domain entities are pure business objects, free from infrastructure concerns
2. **Flexibility**: Can swap databases without changing domain code
3. **Clarity**: Clear separation between "what the business needs" vs "how it's stored"
4. **Evolution**: Database schema can evolve independently from business logic

### Structure

```
internal/
├── domain/
│   ├── entity/          # Domain entities (pure business logic)
│   │   ├── person.go
│   │   ├── school.go
│   │   └── class.go
│   └── repository/      # Repository interfaces (use entities)
│       └── person_repository.go
│
└── repository/
    ├── mapper/          # Convert between entities ↔ models
    │   └── person_mapper.go
    └── sqlite/
        ├── model/       # Persistence models (DB-specific)
        │   └── person.go
        └── person_repository.go  # Implementation (uses models)
```

### Key Principles

#### 1. Repository Interfaces Use Domain Entities
```go
// ✅ GOOD: Interface uses entity
type PersonRepository interface {
    GetByID(ctx context.Context, id uint) (*entity.Person, error)
}

// ❌ BAD: Interface uses persistence model
type PersonRepository interface {
    GetByID(ctx context.Context, id uint) (*model.Person, error)
}
```

#### 2. Repository Implementations Use Persistence Models
```go
// ✅ GOOD: Implementation uses model, maps to entity
func (r *personRepository) GetByID(ctx context.Context, id uint) (*entity.Person, error) {
    var m model.Person
    r.db.First(&m, id)
    return mapper.PersonToEntity(&m), nil
}
```

#### 3. Entities Avoid Deep Relationships
```go
// ✅ GOOD: Use IDs for relationships
type Class struct {
    ID        uint
    TeacherID uint   // ID, not full object
    StudentIDs []uint
}

// ❌ BAD: Deep object relationships
type Class struct {
    Teacher  Person   // Avoids: can cause circular dependencies
    Students []Person // Avoids: loading entire graphs
}
```

### When You Might Skip Mappers

**Only** for very simple CRUD apps where:
- Entities exactly match database schema
- No complex transformations needed
- You're certain you'll never change databases

**But even then**, starting with mappers is safer for future flexibility.

### Benefits Demonstrated

1. **Database Independence**: Change from SQLite to PostgreSQL without touching domain code
2. **Testing**: Mock repository interface with entities, test domain logic in isolation
3. **Performance**: Model can have eager loading, entities remain simple
4. **Schema Evolution**: Add DB columns without breaking domain logic





