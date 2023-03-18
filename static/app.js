function addButtons() {
  var snippets = document.querySelectorAll('.snippet pre');
  var numberOfSnippets = snippets.length;


  for (var i = 0; i < numberOfSnippets; i++) {
    var p = snippets[i].parentElement;
    var b = document.createElement("button");
    b.classList.add('btn-copy')
    b.innerText="Copy";

    b.addEventListener("click", function () {
      this.innerText = 'Copying..';
      code = this.nextSibling.innerText;
      console.log(this.nextSibling);
      navigator.clipboard.writeText(code);
      this.innerText = 'Copied!';
      var that = this;
      setTimeout(function () {
        that.innerText = 'Copy';
      }, 1000)
    });
    p.prepend(b)
  }
}

addButtons();


document.querySelectorAll("pre").forEach(el => el.innerText = el.innerHTML);
