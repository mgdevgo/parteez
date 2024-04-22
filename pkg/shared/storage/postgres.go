package storage

// type CRUDWrapper[TEntity storage.Entity, TID storage.ID] struct {
// 	driver    *pgxpool.Pool
// 	tableName string
// }
//
// func NewCRUDWrapper[TEntity storage.Entity, TID storage.ID](db *pgxpool.Pool, tableName string) CRUDWrapper[TEntity, TID] {
// 	return CRUDWrapper[TEntity, TID]{
// 		driver:    db,
// 		tableName: tableName,
// 	}
// }
//
// func (w CRUDWrapper[TEntity, TID]) Create(entity TEntity) (TID, error) {
// 	const sql = "INSERT INTO %s VALUES (%s) RETURNING id"
//
// 	args, values := make([]any, 0), make([]string, 0)
//
// 	result := w.driver.QueryRow(
// 		context.Background(),
// 		fmt.Sprintf(sql, w.tableName, strings.Join(values, ", ")),
// 		args...,
// 	)
//
// 	var id TID
// 	err := result.Scan(&id)
// 	if err != nil {
// 		return id, err
// 	}
//
// 	return id, nil
// }
//
// func (w CRUDWrapper[TEntity, TID]) Update(id TID, options map[string]any) (TEntity, error) {
// 	//TODO implement me
// 	panic("implement me")
// }
//
// func (w CRUDWrapper[TEntity, TID]) Delete(id TID) error {
// 	//TODO implement me
// 	panic("implement me")
// }
//
// func (w CRUDWrapper[TEntity, TID]) Transaction() (int, error) {
// 	//TODO implement me
// 	//tx, _ := w.driver.Begin(context.Background())
// 	//connection := tx.Conn()
// 	panic("implement me")
// }
