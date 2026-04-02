package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Member struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Points int `json:"points"`
	Tier string `json:"tier"`
	TotalEarned int `json:"total_earned"`
	TotalRedeemed int `json:"total_redeemed"`
	JoinedAt string `json:"joined_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"pennant.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS members(id TEXT PRIMARY KEY,name TEXT NOT NULL,email TEXT DEFAULT '',points INTEGER DEFAULT 0,tier TEXT DEFAULT 'bronze',total_earned INTEGER DEFAULT 0,total_redeemed INTEGER DEFAULT 0,joined_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Member)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO members(id,name,email,points,tier,total_earned,total_redeemed,joined_at,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Email,e.Points,e.Tier,e.TotalEarned,e.TotalRedeemed,e.JoinedAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Member{var e Member;if d.db.QueryRow(`SELECT id,name,email,points,tier,total_earned,total_redeemed,joined_at,created_at FROM members WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Email,&e.Points,&e.Tier,&e.TotalEarned,&e.TotalRedeemed,&e.JoinedAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Member{rows,_:=d.db.Query(`SELECT id,name,email,points,tier,total_earned,total_redeemed,joined_at,created_at FROM members ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Member;for rows.Next(){var e Member;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Points,&e.Tier,&e.TotalEarned,&e.TotalRedeemed,&e.JoinedAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Member)error{_,err:=d.db.Exec(`UPDATE members SET name=?,email=?,points=?,tier=?,total_earned=?,total_redeemed=?,joined_at=? WHERE id=?`,e.Name,e.Email,e.Points,e.Tier,e.TotalEarned,e.TotalRedeemed,e.JoinedAt,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM members WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM members`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Member{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR email LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["tier"];ok&&v!=""{where+=" AND tier=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,email,points,tier,total_earned,total_redeemed,joined_at,created_at FROM members WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Member;for rows.Next(){var e Member;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Points,&e.Tier,&e.TotalEarned,&e.TotalRedeemed,&e.JoinedAt,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
