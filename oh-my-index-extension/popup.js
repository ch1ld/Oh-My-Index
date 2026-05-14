document.getElementById("save").onclick = async () => {

    const tabs = await browser.tabs.query({});

    const data = tabs.map(t => ({
        title: t.title,
        url: t.url
    }));

    fetch("http://127.0.0.1:3737/save", {
        method: "POST",
        body: JSON.stringify(data)
    });

    showToast("Saved successfully!...HonTou?");
};

// 弹窗
const toast = document.getElementById("toast");

function showToast(msg) {
  toast.textContent = msg;
  // 1. 先显示出来
  toast.style.display = "block";
  
  // 2. 微延时，等布局生效再开动画
  setTimeout(() => {
    toast.classList.add("show");
  }, 10);

  // 3. 计时淡出，动画完再隐藏
  setTimeout(() => {
    toast.classList.remove("show");
    // 等透明过渡结束，再彻底关掉 display
    setTimeout(() => {
      toast.style.display = "none";
    }, 300);
  }, 1800);
}