package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`\[\[[A-Za-zěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9|\s]*\]\]`)
	b1 := []byte(`<text xml:space="preserve">&lt;span id=&quot;entry&quot;&gt;von Wieser Friedrich&lt;/span&gt;

	&lt;span class=&quot;PERSON_BORN&quot;&gt;&lt;time datetime=&quot;1851-07-10&quot;&gt;10. července 1851&lt;/time&gt; ve Vídni (Rakousko)&lt;/span&gt;
	&lt;br /&gt;&lt;span class=&quot;PERSON_DIED&quot;&gt;&lt;time datetime=&quot;1926-07-22&quot;&gt;22. července 1926&lt;/time&gt; v Salzburku (Rakousko)&lt;/span&gt;
	&lt;div id=&quot;photo&quot;&gt;[[Soubor:Von Wieser Friedrich 01.jpg|400px]]&lt;/div&gt;
	
	Narodil se a strávil dětství ve Vídni jako syn vysokého úředníka ministerstva války. Studoval práva, byl žákem Carla Mengera, ve své orientaci byl však rovněž ovlivněn čtením ''Úvodu do sociologie'' Herberta Spencera. Později studoval politickou ekonomii v Heidelberku spolu s Eugenem Böhm-Bawerkem, s nímž se spřátelil a jehož sestru si později vzal za manželku. Habilitoval se v roce 1884 a následně byl jmenován docentem na právnické fakultě německé Karlo-Ferdinandovy univerzitě (viz [[Německá univerzita v Praze]]), v roce 1889 pak řádným profesorem. V Praze zůstal do roku 1903, kdy nastoupil po Carlu Mengerovi na vídeňské univerzitě. Spolu s Böhm-Bawerkem formoval první generaci Mengerem založené rakouské ekonomické školy. Jeho žáky byli Ludwig von Mises, Friedrich von Hayek a Joseph A. Schumpeter, kteří pak formovali další generaci této školy. V roce 1917 obdržel baronský titul, stal se členem panské komory rakouského parlamentu a také ministrem obchodu, kterým zůstal do konce první světové války.
	
	Wieser přinesl novou teorii hodnoty, jako první formuloval náklady obětované příležitosti, kategorii marginální užitečnosti (jíž se ale spíše řadí k lausannské škole Léona Walrase a Vilfreda Pareta) a teorii imputace. V posledním čtvrtstoletí se jeho dílo pohybovalo mezi ekonomií, sociologií, politologií a historií. Sociologii, spolu s ekonomií, považoval za disciplínu nutnou k pochopení lidské společnosti a formulaci ekonomické politiky. V díle ''Zákon moci'' (1926) vysvětloval na základě historie vztahy sociálních sil, z nichž ekonomické síly považoval za nejdůležitější pro vývoj společnosti. Striktně odlišoval společenské (sociální) hospodářství (''Gesellschaftliche Wirtschaft'') a socialistickou ekonomii, odmítal kolektivismus, prosazoval metodologický individualismus a ukazoval klíčovou roli jednotlivce pro ekonomickou inovaci. Na jeho ideu vůdcovství v ekonomii později navázal Schumpeter ve své teorii podnikatele a kreativní destrukce.
	
	Český prostor byl významný pro zrození rakouské ekonomické školy, která se nevymezovala přísně vůči sociologii a v některých svých představitelích k ní měla blízko. Zakladatel této školy Menger studoval v Praze, jeho nástupce Wieser v Praze téměř dvacet let přednášel, Böhm-Bawerk pocházel z Brna a v Praze též působil jeho vídeňský student František Čuhel, který se proslavil formulací postoje rakouské školy k uchopení užitku a na jehož odkazu stavěl Mises. Ten, spolu s dalšími představiteli této školy (Hayekem, Machlupem a Schumpeterem, rodákem z Třešti) přenesl po nástupu nacismu její tradici do USA. U nás přetrvala v díle Karla Engliše. Současný ekonom Josef Šíma navázal na tuto tradici počínaje rokem 2005 pořádáním každoročních Conference on Political Economy, v jejichž rámci se vždy koná „wieserovská“ přednáška přednesená některým významným světovým ekonomem, který pak obdrží symbolickou Wieser Memorial Prize.
	
	&lt;span class=&quot;section_title&quot;&gt;Knihy:&lt;/span&gt; ''Die österreichische Schule und die Werth Theorie'' (1891); ''Die Wert Theorie'' (1892); ''Die Theorie der städtischen Grundrente'' (1909); ''Das Wesen und der Hauptinhalt der theoretischen Nationalökonomie'' (1911); ''Theorie der gesellschaftlichen Wirtschaft'' (1914); ''Das geschichtliche Werk der Gewalt'' (1923); ''Die Nationale Steuerleistung und der Landeshaushalt im Königreiche Böhmen: Antwort Auf Die Erwägungen'' (1923); ''Das Gesetz der Macht'' (1926); ''Gessammelte Abhandlungen'' (J. C. B. Mohr, Tübingen 1929; ed. Friedrich A. von Hayek).
	
	&lt;span class=&quot;section_title&quot;&gt;Literatura:&lt;/span&gt; Josef Šíma: Předmluva k Thomas Woods: ''Krach'' (Praha, Dokořán 2010); Milan Sojka: ''Kdo byl kdo – světoví a čeští ekonomové'' (Libri, Praha 2002).
	
	''[[:Kategorie:Aut: Večerník Jiří|Jiří Večerník]]''&lt;br /&gt;
	[[Kategorie:Aut: Večerník Jiří]]
	[[Kategorie:SCSg|Wieser Friedrich von]]</text>`)
	b2 := []byte("Narodil se a strávil dětství ve Vídni jako syn vysokého úředníka ministerstva války. Studoval práva, byl žákem Carla Mengera, ve své orientaci byl však rovněž ovlivněn čtením ''Úvodu do sociologie'' Herberta Spencera. Později studoval politickou ekonomii v Heidelberku spolu s Eugenem Böhm-Bawerkem, s nímž se spřátelil a jehož sestru si později vzal za manželku. Habilitoval se v roce 1884 a následně byl jmenován docentem na právnické fakultě německé Karlo-Ferdinandovy univerzitě (viz [[Německá univerzita v Praze]]), v roce ")
	s := "(viz [[Německá univerzita v Praze]]), v roce"
	fnds := re.FindAllString(s, -1)
	fnd1 := re.FindAll(b1, -1)
	fnd2 := re.FindAll(b2, -1)
	fmt.Println(fnd1, len(fnd2), fnds)
}
