package repository

//func TestNewRepository(t *testing.T) {
//	type args struct {
//		db *sqlx.DB
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Repo
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewRepository() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepo_DeleteUserURL(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx       context.Context
//		shortList []string
//		userID    int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			r.DeleteUserURL(tt.args.ctx, tt.args.shortList, tt.args.userID)
//		})
//	}
//}
//
//func TestRepo_Get(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx      context.Context
//		shortURL string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *dto.Shortening
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			got, err := r.Get(tt.args.ctx, tt.args.shortURL)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Get() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepo_GetByURL(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx       context.Context
//		originURL string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *dto.Shortening
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			got, err := r.GetByURL(tt.args.ctx, tt.args.originURL)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetByURL() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetByURL() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepo_GetListByUser(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx    context.Context
//		userID string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *dto.ShorteningList
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			got, err := r.GetListByUser(tt.args.ctx, tt.args.userID)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetListByUser() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetListByUser() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepo_Put(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx             context.Context
//		shorteningInput dto.Shortening
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *dto.Shortening
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			got, err := r.Put(tt.args.ctx, tt.args.shorteningInput)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Put() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepo_PutList(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx  context.Context
//		list dto.ShorteningList
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			if err := r.PutList(tt.args.ctx, tt.args.list); (err != nil) != tt.wantErr {
//				t.Errorf("PutList() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestRepo_searchURLs(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx     context.Context
//		doneCh  chan struct{}
//		inputCh <-chan string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   <-chan *dto.Shortening
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			if got := r.searchURLs(tt.args.ctx, tt.args.doneCh, tt.args.inputCh); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("searchURLs() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepo_split(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx     context.Context
//		doneCh  chan struct{}
//		inputCh <-chan string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   []<-chan *dto.Shortening
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				db: tt.fields.db,
//			}
//			if got := r.split(tt.args.ctx, tt.args.doneCh, tt.args.inputCh); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("split() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_merge(t *testing.T) {
//	type args struct {
//		doneCh    chan struct{}
//		resultChs []<-chan *dto.Shortening
//	}
//	tests := []struct {
//		name string
//		args args
//		want <-chan *dto.Shortening
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := merge(tt.args.doneCh, tt.args.resultChs...); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("merge() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_prepareList(t *testing.T) {
//	type args struct {
//		doneCh chan struct{}
//		input  []string
//	}
//	tests := []struct {
//		name string
//		args args
//		want <-chan string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := prepareList(tt.args.doneCh, tt.args.input); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("prepareList() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
