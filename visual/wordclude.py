import pymysql
import numpy as np
import jieba
from wordcloud import WordCloud, ImageColorGenerator
from PIL import Image

if __name__ == "__main__":
    db = pymysql.connect("localhost","(づ｡◕‿‿◕｡)づ","123456","bilibili")
    cursor = db.cursor()
    sql = "SELECT title FROM video WHERE pages < 10 ORDER BY view DESC LIMIT 10000"
    cursor.execute(sql)
    results = cursor.fetchall()
    titles = []
    for row in results:
        titles.append(row[0])
    db.close()

    title = " ".join(titles)
    # print(title)

    image = Image.open("assets/mask.png")
    graph = np.array(image)
    # wc = WordCloud(font_path="MonacoYahei.ttf", background_color='white',width=1080,height=1080)
    wc = WordCloud(font_path="assets/MonacoYahei.ttf", background_color='white', mask=graph)
    wc.generate(title)
    wc.to_file("assets/wordcloud.png")
