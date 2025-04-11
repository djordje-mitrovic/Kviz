
package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)


// Ubacuje pitanja u bazu
func InsertQuestions(db *sql.DB) {
	query := `
	INSERT INTO questions (id, pitanja, odgovor1, odgovor2, odgovor3, odgovor4, tacan_odgovor)
	VALUES
	(352, 'Koji je najduži most u svetu?', 'Danyang-Kunshan', 'Venecijanski most', 'Brooklyn', 'Golden Gate', 1),
	(353, 'Koji je najveći morski sisavac?', 'Plavi kit', 'Orka', 'Morski pas', 'Kitovka', 1),
	(354, 'Koji je glavni grad Italije?', 'Rim', 'Milano', 'Napulj', 'Torino', 1),
	(355, 'Koji je najpoznatiji muzejski kompleks u Parizu?', 'Luvr', 'Orsej', 'Panteon', 'Gustav Klimt', 1),
	(356, 'Koja država ima najveći broj ostrva?', 'Švedska', 'Norveška', 'Finska', 'Indonezija', 2),
	(357, 'Koji grad je domaćin Evropskog parlamenta?', 'Strazbur', 'Brisel', 'Amsterdam', 'London', 1),
	(358, 'Koji je najbrži sport na vodi?', 'Jedrenje', 'Veslanje', 'Plivanje', 'Motorski sportovi', 4),
	(359, 'Koja planeta ima najviše vulkana?', 'Venera', 'Mars', 'Jupiter', 'Saturn', 2),
	(360, 'Kako se zvala maskota Olimpijskih igara 1984. u Sarajevu?', 'Vučko', 'Laza', 'Izi', 'Pera', 1),
	(361, 'Koja država je domaćin najviših planina?', 'Nepal', 'Švajcarska', 'Indija', 'Kina', 1),
	(362, 'Koja životinja je najbrža u visini?', 'Soko', 'Jastreb', 'Orao', 'Beli orao', 1),
	(363, 'Koji je glavni grad Egipta?', 'Kairo', 'Alexandria', 'Luxor', 'Giza', 1),
	(364, 'Koje je najdublje jezero u svetu?', 'Bajkalsko', 'Titikaka', 'Kaspijsko more', 'Viktorijino jezero', 1),
	(365, 'Koja je najmlađa planeta u Sunčevom sistemu?', 'Neptun', 'Uran', 'Saturn', 'Jupiter', 1),
	(366, 'Gde se održalo prvo Svetsko prvenstvo u fudbalu?', 'Urugvaj', 'Italija', 'Brazil', 'Argentina', 1),
	(367, 'Koji je najviši građevinski objekat u Sjedinjenim Američkim Državama?', 'Empire State Building', 'Burž Khalifa', 'Petronas Towers', 'One World Trade Center', 4),
	(368, 'Koji kontinent je domaćin najveće džungle?', 'Afrika', 'Azija', 'Srednja Amerika', 'Južna Amerika', 4),
	(369, 'Koji je najveći sisar na kopnu?', 'Afrika slon', 'Indijski slon', 'Nosorog', 'Medved', 1),
	(370, 'Koji je najduži živući organizam?', 'Plavi kit', 'Dugi koralni greben', 'Velika piramida', 'Mamut', 2),
	(371, 'Gde se održalo prvo Svetsko prvenstvo u košarci?', 'Brazil', 'SAD', 'Sovjetski Savez', 'Argentina', 1),
	(372, 'Koja je najviša peščana dina na svetu?', 'Dina Big Daddy', 'Dina Sand', 'Dina Sahara', 'Dina Rub al Khali', 1),
	(373, 'Koja zemlja ima najviše osvojenih Svetskih prvenstava u fudbalu?', 'Brazil', 'Nemačka', 'Argentina', 'Italija', 1),
	(374, 'Koji je najviši narod na svetu?', 'Nizozemci', 'Koreanci', 'Amerikanci', 'Danci', 1),
	(375, 'Koji narod je najniži na svetu?', 'Pigmeji', 'Filipinci', 'Kinezi', 'Indijci', 1),
	(376, 'Koja boja nije na zastavi Francuske?', 'Crvena', 'Plava', 'Bela', 'Crna', 1),
	(377, 'Koja država nema pravougaonu zastavu?', 'Nepal', 'Sjedinjene Američke Države', 'Francuska', 'Brazil', 1),
	(378, 'Koji je najpoznatiji poluotok u Africi?', 'Somalski poluotok', 'Kapski poluotok', 'Sahel', 'Arabijski poluotok', 2),
	(379, 'Koja država je domaćin najduže železničke pruge?', 'Rusija', 'Indija', 'Kina', 'Sjedinjene Američke Države', 1),
	(380, 'Koji grad je poznat po najvećoj količini kiše?', 'Mumbaj', 'Colombo', 'Lima', 'Tucson', 2),
	(381, 'Koji kontinent je domaćin najvećeg broja jezika?', 'Afrika', 'Azija', 'Okeanija', 'Evropa', 1),
	(382, 'Koji je najviši aktivni vulkan na svetu?', 'Etna', 'Kilimandžaro', 'Ojos del Salado', 'Vesuvius', 3),
	(383, 'Koja država ima najviše ostrva?', 'Švedska', 'Indonezija', 'Filipini', 'Norveška', 4),
	(384, 'Koja reka je najduža u Europi?', 'Volga', 'Dunav', 'Tajana', 'Po', 1),
	(385, 'Koji grad je najgušće naseljen na svetu?', 'Džakarta', 'Tokio', 'Mumbaj', 'New York', 2),
	(386, 'Koja je najhladnija pustinja na svetu?', 'Gobi', 'Atakama', 'Sahara', 'Karakum', 1),
	(387, 'Koja zemlja ima najviše planina?', 'Nepal', 'Švajcarska', 'Turska', 'Peru', 1),
	(388, 'Koja je najveća pustinja na svetu?', 'Sahara', 'Gobi', 'Karakum', 'Atakama', 1),
	(389, 'Koji je najveći okean na svetu?', 'Atlantski okean', 'Indijski okean', 'Tihi okean', 'Arktički okean', 3),
	(390, 'Koji je najseverniji glavni grad na svetu?', 'Rejkjavik', 'Oslo', 'Helsinki', 'Sankt Peterburg', 1),
	(391, 'Koja država se prostire na najviše vremenskih zona?', 'Rusija', 'Kanada', 'Sjedinjene Američke Države', 'Brazil', 1),
	(392, 'Koja država je poznata po najdužem obali na svetu?', 'Kanada', 'Australija', 'Rusija', 'Indonezija', 1),
	(393, 'Koji grad je najpoznatiji po gondolama?', 'Venecija', 'Amsterdam', 'Budimpešta', 'Beč', 1),
	(394, 'Koja zemlja ima najveći broj UNESCO spomenika?', 'Italija', 'Kina', 'Francuska', 'Španija', 1),
	(395, 'Koja reka je najduža u Sjedinjenim Američkim Državama?', 'Misouri', 'Misisipi', 'Amazon', 'Colorado', 1),
	(396, 'Koji je najveći grad u Australiji?', 'Sidnej', 'Melburn', 'Brizbej', 'Perte', 1),
	(397, 'Koji je glavni grad Japana?', 'Osaka', 'Tokio', 'Kioto', 'Hirošima', 2),
	(398, 'Koji je najviši vodopad na svetu?', 'Angel', 'Niagara', 'Victoria', 'Iguazu', 1),
	(399, 'Koja država ima najviše jezera na svetu?', 'Kanada', 'Sjedinjene Američke Države', 'Rusija', 'Švedska', 1),
	(400, 'Koji grad je najmanji na svetu?', 'Vatikan', 'Monako', 'Luksemburg', 'Nauru', 1);`
	
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Greška pri ubacivanju podataka:", err)
	}

	fmt.Println("Pitanja su uspešno dodata u bazu!")
}

//
func main() {
	db, err := sql.Open("sqlite3", "quiz.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	InsertQuestions(db)

}

