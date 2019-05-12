# ChatBot

這是一個Keyword Based的Chat Bot函式庫。這個其實跟EverQuest的任務系統有點像，不過強得多。

## Installation

直接以go mod安裝：

```
# In go.mod
require (
	github.com/rayer/chatbot v0.1.0
)

```

## Example

請參考`console_example`目錄

## How to Develop

整個library基本概念就是兩個 : Scenario跟State

### Scenario

「場景」或者說「主題」。比方說，我們可能會提到「關於報告」或者「關於設定」等等的主題，這一個一個的主題就是Scenario。以Mud的講法就是，某個地區。

### State

「狀態」。比方說在這個房子裡面你要探索:

- 你正在一間看起來很大的房間，這房間有往[西邊][南邊]的門跟一個[往下]的樓梯。
- 往下的樓梯
- 你往下走，看到很多廢棄的櫃子堆在那裡。你可以[往回走]，或者[仔細端詳]一下這些櫃子看看有什麼東西
- 端詳櫃子
- 你[看到]一本破舊不堪的筆記本....

這邊一個一個被包起來的就是狀態，而控制這些狀態的就是場景。當然，你可以選擇離開這個空間，每個空間都有不同的狀態。

### Scenario Switch

每個Scenario是可以跳來跳去的，而且，跳過去以後state的狀態會不會被保留也是可以動態決定的。

### Development

比方說，我們想要設計一個能跟你對話，儲存購物清單的bot。那我們可以用一個Scenario叫做

