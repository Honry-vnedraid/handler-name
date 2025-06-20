package main

const (
	GENERATETITLEPROMPT = `Ты получаешь на вход новость: $NEWS1
Достань заголовок новости из поля "text" если он в ней есть. 
Если заголовка нет, то сгенерируй исходя из текста новости.
Заголовок должен быть одним небольшим предложением. 
В качестве ответа укажи только заголовок.

$NEWS1 = %s`
	GETTICKERS = `Ты получаешь на вход новость: $NEWS1
Основываясь на тексте все компании, о которых идет речь в новости. После этого получи тикеры акций/ценных бумаг всех этих компаний, которые сможешь найти.
В качестве источников ищи тикеры как в России, так и в иных странах. 
Если для компании не удалось получить тикер, то скорее всего она Российская и попробуй его езе раз найти с учетом этой информации. 
В качестве ответа укажи ТОЛЬКО тикеры акций/ценных бумаг. Если тикеров не удалось получить, то верни пустую строку. Если таких тикеров больше одного, то перечисли их через запятую с пробелом.

$NEWS1 = %s`
	GETPREDICTIONS = `Ты получаешь на вход новость: $NEWS1 и список тикеров компаний, упомянутых в этой новости.
Представь, что ты биржевой брокер, который только что прочитал эту новость. Считай, что ты владеешь данными тикерами. 
Используя свой многолетний опыт дай прогноз на сколько эта новость может изменить цену отдельно каждого из вышеперечисленных тикеров. 
Дай свою экспертную оценку в виде целого числа в диапазоне от -100 до 100 включительно, где 100 - это очень сильный рост, -100 - очень сильное падение, а 0 - останется без изменений. 
ВАЖНО: Напиши в первой строке через запятую и один пробел оценку для каждого из тикеров в том же порядке, в котором ты их получил. Далее отдельно для КАЖДОГО тикера кратко опиши почему был выставлен такой прогноз, пиши от 3-го лица. Описание для каждого тикера давай в отдельной строке

$NEWS1 = %s

Тикеры: %++v

Пример формата вывода ответа:
20, 60
Продажа непрофильного актива (Дзен) может положительно сказаться на фокусе компании и финансовых показателях, но рынок может воспринять это как признак стратегической неопределенности. Ожидается слабый положительный эффект.
Приобретение Дзен - это стратегическое усиление цифровых активов банка, которое может улучшить его монетизацию и клиентский опыт. Ожидается значительный положительный эффект на котировки.
`

	GETSUMMARY = `Ты получаешь на вход список новостей за день: $LIST
Напиши ТОЛЬКО краткое саммари по самым важным новостям
$LIST = %++v

RESPONSE:

Сбер запустил новую платформу для малого бизнеса, расширив цифровые возможности обслуживания. Тинькофф обновил мобильное приложение с моментальными переводами. Яндекс представил улучшенную голосовую технологию для сервисов. Газпром превысил прогнозы, увеличив добычу газа на 5%% в первом квартале 2025 года.
`
	COMPARENEWS = `Ты получаешь на вход две новости: $NEWS1 и $NEWS2
Твоя задача ответить только ДА или НЕТ
ДА если это одна и та же новость, описанная по разному
НЕТ если это разные новости

$NEWS1 = %++v

$NEWS2 = %++v
	`
)
