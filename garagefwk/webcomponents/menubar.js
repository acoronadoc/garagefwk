
customElements.define("app-menubar", class extends HTMLElement { 

    static observedAttributes = ["sidebarmenu"];
 
    constructor() {
      super();
      this.render();
    }
    
    render() {
        var sidebarmenu = [ ];
        if ( this.getAttribute( "sidebarmenu" ) ) sidebarmenu = JSON.parse( this.getAttribute( "sidebarmenu" ) );

        var menu="";
        for ( var i=0; i<sidebarmenu.length; i++ ) {
            if ( sidebarmenu[i].Title ) menu += "<div class='title'>" + sidebarmenu[i].Title + "</div>";

            if ( sidebarmenu[i].Childs )
                for ( var n=0; n<sidebarmenu[i].Childs.length; n++ ) {
                    menu += "<a href='"+ sidebarmenu[i].Childs[n].Href +"'><div class='optionmenu'>"+ sidebarmenu[i].Childs[n].Title +"</div></a>";
                }
        }

        this.innerHTML = `
            <div id="headerside">
                <a href="#" id="ham-main-menu"></a>
            </div>
            <div id="header-menu">
                <div id="header-bg"></div>
                <div id="header-right">${menu}</div>
            </div>
        `;

        var el = this;
        this.querySelector("#ham-main-menu").addEventListener( "click", function(event) {
            el.querySelector("#header-menu").classList.toggle('open'); 
            return false;
        } );

        this.querySelector("#header-bg").addEventListener( "click", function(event) {
            el.querySelector("#header-menu").classList.toggle('open'); 
            return false;
        } );

    }
});