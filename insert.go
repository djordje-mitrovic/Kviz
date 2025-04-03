
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
	(1, 'Koji je glavni grad Francuske?', 'Madrid', 'London', 'Pariz', 'Berlin', 3),
	(2, 'Kad je bio Kosovski boj?', '1389', '1375', '1412', '1403', 1),
	(3, 'Koji je hemijski simbol za vodonik?', 'O', 'H', 'C', 'N', 2),
	(4, 'Koja je najveća pustinja na svetu?', 'Sahara', 'Gobi', 'Antarktička', 'Kalahari', 3),
	(5, 'Ko je naslikao Mona Lizu?', 'Van Gog', 'Mikelanđelo', 'Leonardo da Vinči', 'Renoar', 3),
	(6, 'Koja država je osvojila FIFA Svetsko prvenstvo 2018. godine?', 'Brazil', 'Francuska', 'Nemačka', 'Argentina', 2),
	(7, 'Koji je hemijski simbol za zlato?', 'Ag', 'Au', 'Pb', 'Zn', 2),
	(8, 'Ko je napisao roman "Zločin i kazna"?', 'Tolstoj', 'Dostojevski', 'Gogolj', 'Čehov', 2),
	(9, 'Kako se zvao prvi čovek koji je kročio na Mesec?', 'Jurij Gagarin', 'Nil Armstrong', 'Buzz Aldrin', 'Glen', 2),
	(10, 'Koji okean je najveći na svetu?', 'Atlantski', 'Indijski', 'Tihi', 'Severni ledeni', 3),
	(11, 'Koji sport je najpopularniji u Brazilu?', 'Ragbi', 'Košarka', 'Fudbal', 'Tenis', 3),
	(12, 'Koji glumac igra glavnu ulogu u filmu "Titanik"?', 'Leonardo Dikaprio', 'Tom Henks', 'Bred Pit', 'Džoni Dep', 1),
	(13, 'Koji je najveći kontinent na svetu?', 'Evropa', 'Afrika', 'Azija', 'Severna Amerika', 3),
	(14, 'Ko je napisao "Božanstvenu komediju"?', 'Homer', 'Vergilije', 'Dante Aligijeri', 'Šekspir', 3),
	(15, 'Koja je najduža reka na svetu?', 'Nil', 'Amazona', 'Misisipi', 'Jangce', 2),
	(16, 'Ko je osnivač Microsofta?', 'Stiv Džobs', 'Bil Gejts', 'Mark Zakerberg', 'Elon Mask', 2),
	(17, 'Koji je simbol rimskog boga Marsa?', 'Koplje', 'Luk', 'Čekić', 'Mač', 1),
	(18, 'Koja je nacionalna valuta Japana?', 'Dolar', 'Jen', 'Juan', 'Evro', 2),
	(19, 'Koja planina je najviša na svetu?', 'Mont Everest', 'K2', 'Kilimandžaro', 'Alpi', 1),
	(20, 'Koji pisac je napisao "Ana Karenjina"?', 'Tolstoj', 'Dostojevski', 'Gogolj', 'Puškin', 1),
	(21, 'Koji je prvi film u istoriji kinematografije?', 'Građanin Kejn', 'Dolazak voza', 'Metropolis', 'Prohujalo s vihorom', 2),
	(22, 'Koja država je domaćin Olimpijskih igara 2024. godine?', 'SAD', 'Kina', 'Francuska', 'Brazil', 3),
	(23, 'Koji je zvanični jezik u Brazilu?', 'Španski', 'Portugalski', 'Engleski', 'Francuski', 2),
	(24, 'Koji organ u ljudskom telu filtrira krv?', 'Jetra', 'Bubrezi', 'Pluća', 'Slezina', 2),
	(25, 'Kako se zove najviši toranj u svetu?', 'Ajfelov toranj', 'Burž Kalifa', 'Empire State Building', 'CN Tower', 2),
	(26, 'Koji klub je osvojio UEFA Ligu šampiona 2022. godine?', 'Barselona', 'Real Madrid', 'Bajern Minhen', 'Liverpul', 2),
	(27, 'Koji kompozitor je napisao "Mesečevu sonatu"?', 'Bach', 'Beethoven', 'Mozart', 'Čajkovski', 2),
	(28, 'Koja država ima najviše piramida na svetu?', 'Egipat', 'Meksiko', 'Sudan', 'Peru', 3),
	(29, 'Koji gas je najzastupljeniji u Zemljinoj atmosferi?', 'Oksigen', 'Azot', 'Ugljen-dioksid', 'Argon', 2),
	(30, 'Koji pisac je autor romana "1984"?', 'Džordž Orvel', 'Oldos Haksli', 'Ričard Bah', 'Rej Bredberi', 1),
	(31, 'Kako se zove najpoznatija slika Edvarda Munka?', 'Vrisak', 'Mona Liza', 'Zvezdana noć', 'Dama sa hermelinom', 1),
	(32, 'Koji je najviši vodopad na svetu?', 'Viktorijini vodopadi', 'Niagarini vodopadi', 'Anđeoski vodopad', 'Iguazu', 3),
	(33, 'Ko je režirao film "Interstellar"?', 'Stiven Spilberg', 'Martin Skorseze', 'Kventin Tarantino', 'Kristofer Nolan', 4),
	(34, 'Koji metal je teži: gvožđe ili aluminijum?', 'Aluminijum', 'Gvožđe', 'Oba su ista', 'Zavisi od oblika', 2),
	(35, 'Koja država se prostire na dva kontinenta?', 'Egipat', 'Turska', 'Rusija', 'Kanada', 3),
	(36, 'Koja je najpoznatija opera Verdija?', 'La Traviata', 'Seviljski berberin', 'Figarova ženidba', 'Toska', 1),
	(37, 'Ko je otkrio Ameriku?', 'Vasko de Gama', 'Kristofer Kolumbo', 'Magelan', 'Džejms Kuk', 2),
	(38, 'Koja je osnovna jedinica mere za masu?', 'Kilogram', 'Litra', 'Newton', 'Gram', 1),
	(39, 'Koja država ima najviše osvojenih zlatnih medalja na Olimpijskim igrama?', 'Kina', 'Rusija', 'SAD', 'Nemačka', 3),
	(40, 'Ko je napisao roman "Ponos i predrasude"?', 'Džejn Ostin', 'Šarlot Bronte', 'Emili Bronte', 'Luiza Mej Alkot', 1),
	(41, 'Koji je najveći cvet na svetu?', 'Rafflezija', 'Sunovrat', 'Lotos', 'Orhideja', 1),
	(42, 'Koji je glavni sastojak sušija?', 'Riba', 'Pirinač', 'Morska trava', 'Soja sos', 2),
	(43, 'Ko je bio prvi car Rimskog carstva?', 'Julije Cezar', 'Avgust', 'Neron', 'Kaligula', 2),
	(44, 'Kako se zove mitološko biće sa telom konja i ljudskim torzom?', 'Minotaur', 'Kentaur', 'Sfinga', 'Pegaz', 2),
	(45, 'Koji sport koristi termin "love" za rezultat nula?', 'Fudbal', 'Košarka', 'Tenis', 'Golf', 3),
	(46, 'Koji je glavni sastojak piva?', 'Voda', 'Ječam', 'Hmelj', 'Svi navedeni', 4),
	(47, 'Koji naučnik je postavio teoriju relativnosti?', 'Njutn', 'Ajnštajn', 'Tesla', 'Boh', 2),
	(48, 'Koji grad je poznat kao "vetroviti grad"?', 'Njujork', 'Čikago', 'London', 'Pariz', 2),
	(49, 'Koja životinja simbolizuje mudrost?', 'Lav', 'Sova', 'Zmija', 'Slon', 2),
	(50, 'Ko je osnivač Apple-a?', 'Bil Gejts', 'Stiv Džobs', 'Mark Zakerberg', 'Elon Mask', 2),
	(51, 'Koji grad je poznat po kanalu Veliki kanal?', 'Venecija', 'Amsterdam', 'Pariz', 'Sankt Peterburg', 1),
	(52, 'Ko je napisao roman "Rat i mir"?', 'Tolstoj', 'Dostojevski', 'Gogolj', 'Čehov', 1),
	(53, 'Kako se zove glavni lik u romanu "Don Kihot"?', 'Sančo Pansa', 'Don Kihot', 'Hamlet', 'Raskoljnikov', 2),
	(54, 'Koji je najveći stadion na svetu po kapacitetu?', 'Kamp Nou', 'Rungrado May Day', 'Marakana', 'Vembli', 2),
	(55, 'Koja zemlja je izgradila prvi metro?', 'Francuska', 'SAD', 'Velika Britanija', 'Rusija', 3),
	(56, 'Koji gas koristimo za disanje?', 'Azot', 'Oksigen', 'Ugljen-dioksid', 'Vodonik', 2),
	(57, 'Koja ptica može da leti unazad?', 'Soko', 'Kolibri', 'Orao', 'Svraka', 2),
	(58, 'Kako se zove najpoznatiji toranj u Pizi?', 'Krivi toranj', 'Ajfelov toranj', 'Torre de Belém', 'Big Ben', 1),
	(59, 'Ko je bio prvi ruski car?', 'Ivan Grozni', 'Petar Veliki', 'Nikola II', 'Aleksandar I', 1),
	(60, 'Koja država ima najviše stanovnika na svetu?', 'SAD', 'Indija', 'Kina', 'Rusija', 3),
	(61, 'Koji sport se igra na Vimbldonu?', 'Fudbal', 'Košarka', 'Tenis', 'Golf', 3),
	(62, 'Kako se zove čuveni američki filmski festival?', 'Kanski festival', 'Oskari', 'Sundance', 'Venecijanski festival', 3),
	(63, 'Koja reka protiče kroz Egipat?', 'Nil', 'Amazona', 'Misisipi', 'Dunav', 1),
	(64, 'Koji kontinent nema stalno naseljeno stanovništvo?', 'Afrika', 'Antarktik', 'Australija', 'Evropa', 2),
	(65, 'Ko je bio prvi rimski imperator?', 'Julije Cezar', 'Neron', 'August', 'Kaligula', 3),
	(66, 'Koji instrument ima dirke i pedale?', 'Gitara', 'Violina', 'Klavir', 'Harmonika', 3),
	(67, 'Koja država ima najdužu obalu na svetu?', 'SAD', 'Kanada', 'Australija', 'Rusija', 2),
	(68, 'Koji planet je najbliži Suncu?', 'Mars', 'Venera', 'Merkur', 'Jupiter', 3),
	(69, 'Kako se zove najviši vulkan na svetu?', 'Etna', 'Kilimandžaro', 'Ojos del Salado', 'Fudži', 3),
	(70, 'Ko je napisao "Hamleta"?', 'Dante', 'Šekspir', 'Molijer', 'Gogolj', 2),
	(71, 'Koja je najhladnija planeta Sunčevog sistema?', 'Saturn', 'Neptun', 'Mars', 'Uran', 2),
	(72, 'Koji glumac je poznat po ulozi Džejmsa Bonda?', 'Šon Koneri', 'Džoni Dep', 'Leonardo Dikaprio', 'Bred Pit', 1),
	(73, 'Koji grad se naziva "Grad svetlosti"?', 'London', 'Pariz', 'Njujork', 'Rim', 2),
	(74, 'Koji element ima hemijski simbol Fe?', 'Bakar', 'Gvožđe', 'Cink', 'Srebro', 2),
	(75, 'Koji je najveći sisar na svetu?', 'Slon', 'Plavi kit', 'Nosorog', 'Gorila', 2),
	(76, 'Koja boja se dobija mešanjem plave i žute?', 'Zelena', 'Narandžasta', 'Ljubičasta', 'Braon', 1),
	(77, 'Kako se zove najveće more na svetu?', 'Crno more', 'Sredozemno more', 'Južno kinesko more', 'Kaspijsko more', 4),
	(78, 'Ko je otkrio zakon gravitacije?', 'Ajnštajn', 'Njutn', 'Tesla', 'Galileo', 2),
	(79, 'Kako se zove glavni lik u "Harry Potter" serijalu?', 'Ron', 'Harry', 'Hermiona', 'Dambldor', 2),
	(80, 'Koji sport se igra sa palicom i lopticom na ledu?', 'Hokej', 'Bejzbol', 'Ragbi', 'Fudbal', 1),
	(81, 'Koja je najviša zgrada u Americi?', 'Empire State Building', 'Burž Kalifa', 'One World Trade Center', 'Willis Tower', 3),
	(82, 'Koji brod je potonuo 1912. godine?', 'Titanik', 'Lusitanija', 'Britanik', 'Endeavour', 1),
	(83, 'Ko je osnivač Facebook-a?', 'Elon Musk', 'Bil Gejts', 'Mark Zakerberg', 'Stiv Džobs', 3),
	(84, 'Koji je glavni grad Kanade?', 'Toronto', 'Montreal', 'Otava', 'Vankuver', 3),
	(85, 'Koja životinja može da spava stojeći?', 'Slon', 'Konj', 'Lav', 'Panda', 2),
	(86, 'Koja je najveća planeta Sunčevog sistema?', 'Mars', 'Saturn', 'Jupiter', 'Neptun', 3),
	(87, 'Koja ptica ne može da leti?', 'Orao', 'Pingvin', 'Soko', 'Galeb', 2),
	(88, 'Koja država je domaćin karnevala u Riju?', 'Argentina', 'Meksiko', 'Brazil', 'Kolumbija', 3),
	(89, 'Koji je najtopliji mesec u godini u Evropi?', 'Januar', 'Jul', 'Septembar', 'Mart', 2),
	(90, 'Koji je prvi film iz serijala "Ratovi zvezda"?', 'Epizoda IV', 'Epizoda I', 'Epizoda V', 'Epizoda VI', 1),
	(91, 'Koji je najstariji univerzitet u Evropi?', 'Kembridž', 'Oksford', 'Bolonja', 'Sorbon', 3),
	(92, 'Kako se zove grčki bog rata?', 'Had', 'Zeus', 'Ares', 'Hermes', 3),
	(93, 'Ko je napisao roman "Derviš i smrt"?', 'Ivo Andrić', 'Mesa Selimović', 'Dobrica Ćosić', 'Branko Ćopić', 2),
	(94, 'Koji metal se koristi za izradu električnih kablova?', 'Željezo', 'Bakar', 'Aluminijum', 'Olovо', 2),
	(95, 'Koja rijeka protiče kroz London?', 'Dunav', 'Temza', 'Rajna', 'Seina', 2),
	(96, 'Koja je glavna boja na kineskoj zastavi?', 'Crvena', 'Plava', 'Bela', 'Žuta', 1),
	(97, 'Ko je bio prvi čovek u svemiru?', 'Nil Armstrong', 'Jurij Gagarin', 'Buzz Aldrin', 'Majkl Kolins', 2),
	(98, 'Kako se zove festival filma u Kanu?', 'Venecijanski festival', 'Berlinale', 'Kanski festival', 'Oskari', 3),
	(99, 'Ko je autor romana "Braća Karamazovi"?', 'Tolstoj', 'Dostojevski', 'Gogolj', 'Turgenjev', 2),
	(100, 'Ko je bio prvi predsednik SAD?', 'Džon Adams', 'Tomas Džeferson', 'Džordž Vašington', 'Džejms Medison', 3);
	`
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

