---
title: "ASCII码表"
date: 2023-02-23
draft: false
tags : [                    # 文章所属标签
    "Linux",
]
categories : [              # 文章所属标签
    "技术",
]
---

ASCII码，使用7位二进制数，表示128个标准ASCII字符，使用8位二进制数，表示256 个标准及扩展ASCII字符；

ASCII编码字符分类：

控制字符：0～32、127表示，共33个，如CR（回车）、LF（换行）、FF（换页）、BS（退格）、DEL（删除）、Space（空格）等。

特殊符号：33-47表示，如+（加）、-（减）、*（乘）、/（除）、！（感叹号）。

数字：48～57表示，0-9阿拉伯数字。

字母：65～90为26个大写英文字母，97～122号为26个小写英文字母。


|DEC|OCT|HEX|BIN|缩写/符号|HTML实体|描述|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|0|000|00|00000000|NUL|�|Null char (空字符)|
|1|001|01|00000001|SOH|�|Start of Heading (标题开始)|
|2|002|02|00000010|STX|�|Start of Text (正文开始)|
|3|003|03|00000011|ETX|�|End of Text (正文结束)|
|4|004|04|00000100|EOT|�|End of Transmission (传输结束)|
|5|005|05|00000101|ENQ|�|Enquiry (请求)|
|6|006|06|00000110|ACK|�|Acknowledgment (收到通知)|
|7|007|07|00000111|BEL|�|Bell (响铃)|
|8|010|08|00001000|BS|�|Back Space (退格)|
|9|011|09|00001001|HT||Horizontal Tab (水平制表符)|
|10|012|0A|00001010|LF||Line Feed (换行键)|
|11|013|0B|00001011|VT|�|Vertical Tab (垂直制表符)|
|12|014|0C|00001100|FF||Form Feed (换页键)|
|13|015|0D|00001101|CR||Carriage Return (回车键)|
|14|016|0E|00001110|SO|�|Shift Out / X-On (不用切换)|
|15|017|0F|00001111|SI|�|Shift In / X-Off (启用切换)|
|16|020|10|00010000|DLE|�|Data Line Escape (数据链路转义)|
|17|021|11|00010001|DC1|�|Device Control 1 (设备控制1)|
|18|022|12|00010010|DC2|�|Device Control 2 (设备控制2)|
|19|023|13|00010011|DC3|�|Device Control 3 (设备控制3)|
|20|024|14|00010100|DC4|�|Device Control 4 (设备控制4)|
|21|025|15|00010101|NAK|�|Negative Acknowledgement (拒绝接收)|
|22|026|16|00010110|SYN|�|Synchronous Idle (同步空闲)|
|23|027|17|00010111|ETB|�|End of Transmit Block (传输块结束)|
|24|030|18|00011000|CAN|�|Cancel (取消)|
|25|031|19|00011001|EM|�|End of Medium (介质中断)|
|26|032|1A|00011010|SUB|�|Substitute (替补)|
|27|033|1B|00011011|ESC|�|Escape (溢出)|
|28|034|1C|00011100|FS|�|File Separator (文件分割符)|
|29|035|1D|00011101|GS|�|Group Separator (分组符)|
|30|036|1E|00011110|RS|�|Record Separator (记录分离符)|
|31|037|1F|00011111|US|�|Unit Separator (单元分隔符)|
|32|040|20|00100000|||Space (空格)|
|33|041|21|00100001|!|!|Exclamation mark|
|34|042|22|00100010|"|"|Double quotes|
|35|043|23|00100011|#|#|Number|
|36|044|24|00100100|$|$|Dollar|
|37|045|25|00100101|%|%|Procenttecken|
|38|046|26|00100110|&|&|Ampersand|
|39|047|27|00100111|’|'|Single quote|
|40|050|28|00101000|(|(|Open parenthesis\
|41|051|29|00101001|)|)|Close parenthesis|
|42|052|2A|00101010|*|*|Asterisk|
|43|053|2B|00101011|+|+|Plus|
|44|054|2C|00101100|,|,|Comma|
|45|055|2D|00101101|-|-|Hyphen|
|46|056|2E|00101110|.|.|Period, dot or full stop|
|47|057|2F|00101111|/|/|Slash or divide|
|48|060|30|00110000|0|0|Zero|
|49|061|31|00110001|1|1|One|
|50|062|32|00110010|2|2|Two|
|51|063|33|00110011|3|3|Three|
|52|064|34|00110100|4|4|Four|
|53|065|35|00110101|5|5|Five|
|54|066|36|00110110|6|6|Six|
|55|067|37|00110111|7|7|Seven|
|56|070|38|00111000|8|8|Eight|
|57|071|39|00111001|9|9|Nine|
|58|072|3A|00111010|:|:|Colon|
|59|073|3B|00111011|;|;|Semicolon|
|60|074|3C|00111100|<|<|Less than|
|61|075|3D|00111101|=|=|Equals|
|62|076|3E|00111110|>|>|Greater than|
|63|077|3F|00111111|?|?|Question mark|
|64|100|40|01000000|@|@|At symbol|
|65|101|41|01000001|A|A|Uppercase A|
|66|102|42|01000010|B|B|Uppercase B|
|67|103|43|01000011|C|C|Uppercase C|
|68|104|44|01000100|D|D|Uppercase D|
|69|105|45|01000101|E|E|Uppercase E|
|70|106|46|01000110|F|F|Uppercase F|
|71|107|47|01000111|G|G|Uppercase G|
|72|110|48|01001000|H|H|Uppercase H|
|73|111|49|01001001|I|I|Uppercase I|
|74|112|4A|01001010|J|J|Uppercase J|
|75|113|4B|01001011|K|K|Uppercase K|
|76|114|4C|01001100|L|L|Uppercase L|
|77|115|4D|01001101|M|M|Uppercase M|
|78|116|4E|01001110|N|N|Uppercase N|
|79|117|4F|01001111|O|O|Uppercase O|
|80|120|50|01010000|P|P|Uppercase P|
|81|121|51|01010001|Q|Q|Uppercase Q|
|82|122|52|01010010|R|R|Uppercase R|
|83|123|53|01010011|S|S|Uppercase S|
|84|124|54|01010100|T|T|Uppercase T|
|85|125|55|01010101|U|U|Uppercase U|
|86|126|56|01010110|V|V|Uppercase V|
|87|127|57|01010111|W|W|Uppercase W|
|88|130|58|01011000|X|X|Uppercase X|
|89|131|59|01011001|Y|Y|Uppercase Y|
|90|132|5A|01011010|Z|Z|Uppercase Z|
|91|133|5B|01011011|[|[|Opening bracket|
|92|134|5C|01011100|\\|\\|Backslash|
|93|135|5D|01011101|]|]|Closing bracket|
|94|136|5E|01011110|^|^|Caret - circumflex|
|95|137|5F|01011111|_|_|Underscore|
|96|140|60|01100000|\`|\`|Grave accent|
|97|141|61|01100001|a|a|Lowercase a|
|98|142|62|01100010|b|b|Lowercase b|
|99|143|63|01100011|c|c|Lowercase c|
|100|144|64|01100100|d|d|Lowercase d|
|101|145|65|01100101|e|e|Lowercase e|
|102|146|66|01100110|f|f|Lowercase f|
|103|147|67|01100111|g|g|Lowercase g|
|104|150|68|01101000|h|h|Lowercase h|
|105|151|69|01101001|i|i|Lowercase i|
|106|152|6A|01101010|j|j|Lowercase j|
|107|153|6B|01101011|k|k|Lowercase k|
|108|154|6C|01101100|l|l|Lowercase l|
|109|155|6D|01101101|m|m|Lowercase m|
|110|156|6E|01101110|n|n|Lowercase n|
|111|157|6F|01101111|o|o|Lowercase o|
|112|160|70|01110000|p|p|Lowercase p|
|113|161|71|01110001|q|q|Lowercase q|
|114|162|72|01110010|r|r|Lowercase r|
|115|163|73|01110011|s|s|Lowercase s|
|116|164|74|01110100|t|t|Lowercase t|
|117|165|75|01110101|u|u|Lowercase u|
|118|166|76|01110110|v|v|Lowercase v|
|119|167|77|01110111|w|w|Lowercase w|
|120|170|78|01111000|x|x|Lowercase x|
|121|171|79|01111001|y|y|Lowercase y|
|122|172|7A|01111010|z|z|Lowercase z|
|123|173|7B|01111011|{|{|Opening brace|
|124|174|7C|01111100|||||Vertical bar|
|125|175|7D|01111101|}|}|Closing brace|
|126|176|7E|01111110|~|~|Equivalency sign (tilde)|
|127|177|7F|01111111||�|Delete|

————————————————
版权声明：本文为CSDN博主「火腿肠」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/qq_39511050/article/details/126809454