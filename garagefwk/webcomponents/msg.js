
customElements.define("app-msg", class extends HTMLElement { 

    static observedAttributes = ["sidebarmenu"];
 
    constructor() {
      super();
    }

    connectedCallback() {
        this.render();
    }
     
    attributeChangedCallback(name, oldValue, newValue) {
        this.render();
    }    

    render() {
        this.innerHTML=`
            <a href='#'></a>
            ${this.getAttribute('msg')}
        `;

        var el=this;
        this.querySelector("a").addEventListener("click", function(event) {
            el.className += " removing";
            setTimeout( function() { el.remove(); }, 500 );
        });
    }

});
