<template>
  <div>
    <aside>
      <nav>
        <div class="accordion" id="accordionPanelsStayOpenExample">
          <div class="accordion-item">
            <h2 class="accordion-header" id="panelsStayOpen-headingOne">
              <button
                class="accordion-button"
                type="button"
                data-bs-toggle="collapse"
                data-bs-target="#panelsStayOpen-collapseOne"
                aria-expanded="true"
                aria-controls="panelsStayOpen-collapseOne"
              >
                Chatrooms
              </button>
            </h2>
            <div
              id="panelsStayOpen-collapseOne"
              class="accordion-collapse collapse show"
              aria-labelledby="panelsStayOpen-headingOne"
            >
              <div class="accordion-body">
                <button type="button" class="btn btn-primary">
                  create chatroom
                </button>
                <ul>
                  <li><router-link to="/chatroom/1">chat 1</router-link></li>
                  <li><router-link to="/chatroom/2">Friends</router-link></li>
                  <li><router-link to="/chatroom/54">Stocks</router-link></li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </nav>
    </aside>

    <section>
      <router-view></router-view>
    </section>
  </div>
</template>

<style scoped>
aside,
section {
  display: inline-block;
  vertical-align: top;
  width: 30%;
}

section {
  width: 65%;
}
</style>