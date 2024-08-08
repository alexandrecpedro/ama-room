<div align = "center">
  <h1>AMA (Ask me anything) Room</h1>
</div>
<br>

<div align = 'center' justify-content = 'space-around' >
  <img width="1604" alt="Rocketseat AMA Room - Desktop" src="./project/screens/screen1.png">
</div>
<br>
<br>

<p></p>

<p align="center">
 <a href="#theproject">The Project</a> •
 <a href="#target">Target</a> •
 <a href="#technologies">Technologies</a> •
 <a href="#route">Route</a> •
 <a href="#howtouse">How to Use</a>
</p>
<br>

<div id="theproject">
<h2> 📓 The Project </h2>
<p> AMA (Ask me anything) Room is a public mentoring app where developers send questions and the community votes on the coolest questions in real time via WebSockets</p>
</div>

<div id="target">
<h2> 💡 Target </h2>
Development of a AMA (Ask me anything) Room system, a public mentoring app where developers send questions and the community votes on the coolest questions in real time via WebSockets at Tech Week Go+React, from Rocketseat
</div>
<br>

<div id="technologies">
<h2> 🛠 Technologies </h2>
The following tools were used in building the project:<br><br>

|                      Type                       |       Tools       |              References               |
|:-----------------------------------------------:|:-----------------:|:-------------------------------------:|
|                       IDE                       |    VISUAL CODE    |    https://code.visualstudio.com/     |
|              Design Interface Tool              | FIGMA (Prototype - UX/UI) |      https://www.figma.com/   |
|         Programming Language (Frontend)         |       REACT       |       https://reactjs.org/            |
|         Programming Language (Frontend)         |    TYPESCRIPT     |  https://www.typescriptlang.org/      |
|      Utility-first CSS Framework (Frontend)     |   TAILWIND CSS    |     https://tailwindcss.com/          |
|    Tool for transforming CSS with JavaScript    |     POST CSS      |       https://postcss.org/            |
|         Graphic components (Frontend)           |  PHOSPHOR ICONS   |    https://phosphoricons.com/         |
|             UI Components for React             |     RADIX-UI      |     https://www.radix-ui.com/         |
|     Tool to build frontend faster (Frontend)    |      VITE.JS      |        https://vitejs.dev/            |
|  Promise based HTTP client - browser & Node.js  |       AXIOS       |      https://axios-http.com/          |
|       API and backend services (Backend)        |        GO         |     https://go.dev/                   |
|               Database (Backend)                |    POSTGRESQL     |   https://www.postgresql.org/         |
| Open source API development ecosystem (Testing) |      POSTMAN      |       https://www.postman.com/        |
</div>
<br>

<div align = 'center'>
  <h3>Backend | API</h3>
  <img height =' 100px ' left='80px' src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/go/go-original-wordmark.svg" />
  <!-- <img height =' 100px ' left='80px' src="https://jwt.io/img/logo-asset.svg" /> -->
  <br>
  <h3>Testing</h3>
  <img width =' 100px ' src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/postman/postman-original.svg" />
  <br>
  <h3>Database</h3>
  <img height =' 100px ' src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/postgresql/postgresql-original.svg" />
  <br>
  <h3>IDE</h3>
  <img height =' 100px ' src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/vscode/vscode-original.svg" />
  <br>
  <!-- <h3>UX/UI</h3>
  <img height =' 100px ' src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/figma/figma-original.svg" />
  <img height =' 100px ' left='80px' src="./project/logo/phosphor-icons_logo.png"/>
  <br> -->
  <h3>Frontend</h3>
  <img width =' 100px ' src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg" />
  <img width =' 100px ' left='80px' src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/typescript/typescript-original.svg" />
  <!-- <img height =' 100px ' left='80px' src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nextjs/nextjs-original-wordmark.svg"/> -->
  <br>
  <img height =' 50px ' src="./project/logo/tailwind-css_logo.svg" />
  <br>
</div>
<br>

<div id="route">
<h2> 🔎 Route </h2>
  <ol>
    <li &nbsp;>Part 1 - Go + React | Class 01</li>
    <br>
    <li &nbsp;>Part 2 - Go + React | Class 02</li>
    <br>
    <li &nbsp;>Part 3 - Go + React | Class 03</li>
    <br>
    <li &nbsp;>Part 4 - Go + React | Class 04</li>
    <br>
  </ol>
</div>

<div id="route">
<h2> 🔎 Route </h2>
  <ol>
    <li &nbsp;>Part 1 - Go<br/>
      <ul &nbsp;>
        <li &nbsp;>Build the project prototype: https://www.figma.com/file/zYhLAnJ0IPZjN4FsZ0GX8r/Habits-(i)-(Community)?node-id=6%3A343&t=eYeBEtKtUvGaZ9zz-0</li>
        <li>Install VS Code (IDE)</li>
        <li>Install VS Code extensions: NodeJS, Prisma, Tailwind CSS IntelliSense, PostCSS Language Support, Symbols, Fluent Icons, DotENV, Go</li>
        <li &nbsp;><b>Backend project</b>
          <ul>
            <li>Create a new project: mkdir backend</li>
            <li>Enter backend project: cd backend</li>
            <li>Install Go and start: go mod init GITHUB_REPO_PATH (without https://)</li>
            <li>Create command directory: ./cmd</li>
            <li>Create directory for usage for this project or executable packages: ./internal</li>
            <li>Database
              <ol>
                <li>Settings: compose.yml</li>
                <li>TERN services
                  <ul>
                    <li>Lib - TERN: https://github.com/jackc/tern</li>
                    <li>Install: go install github.com/jackc/tern/v2@latest</li>
                    <li>Create directory for migrations: mkdir internal/store/pgstore</li>
                    <li>Init: tern init ./internal/store/pgstore/migrations</li>
                    <li>Set TERN configuration</li>
                    <li>Create migration: tern new --migrations ./internal/store/pgstore/migrations MIGRATION_NAME</li>
                  </ul>
                </li>
                <li>Create command to run migrations: ./cmd/tools/terndotenv/main.go</li>
                <li>Install godotenv lib: go install github.com/joho/godotenv/cmd/godotenv@latest</li>
                <li>Pass this dependency to go.mod: go mod tidy</li>
                <li>Execute migrations: go run cmd/tools/terndotenv/main.go</li>
              </ol>
            </li>
            <li>Representing database tables on Go
              <ol>
                <li>Package - SQLC: https://sqlc.dev/</li>
                <li>Install: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest</li>
                <li>Set config file: ./internal/store/pgstore/sqlc.yaml</li>
                <li>Create queries
                  <ul>
                    <li>Create folder: mkdir ./internal/store/pgstore/queries</li>
                    <li>Set queries: touch ./internal/store/pgstore/queries/queries.sql</li>
                  </ul>
                </li>
                <li>Generate code: sqlc generate -f ./internal/store/pgstore/sqlc.yaml</li>
                <li>Run: go mod tidy</li>
              </ol>
            </li>
            <li>Go directive tool: gen.go
              <ul>
                <li>Simplify usage of run migration and generate sqlc code</li>
                <li>Run directive: go generate ./...</li>
              </ul>
            </li>
            <li>Write API basics: ./internal/api/api.go</li>
            <li>Create main function: ./cmd/wsrs/main.go</li>
            <li>Start server: go run ./cmd/wsrs/main.go</li>
          </ul>
        </li>
      </ul>
    </li>
    <br>
    <li &nbsp;>Part 2 - Go<br/>
      <ul &nbsp;>
        <li &nbsp;><b>Backend project</b>
          <ul>
            <li>Install libs
              <ol>
                <li>Gorilla Websocket: go get github.com/gorilla/websocket</li>
                <li></li>
                <li></li>
              </ol>
            </li>
            <li>API configuration: ./internal/api/api.go
              <ol>
                <li>Set middlewares</li>
                <li>CORS
                  <ul>
                    <li>Install Chi CORS lib: go get github.com/go-chi/cors</li>
                    <li>Apply basic CORS</li>
                  </ul>
                </li>
                <li>Set routes</li>
                <li>Implementing API methods from routes</li>
              </ol>
            </li>
          </ul>
        </li>
        <li &nbsp;><b>Testing</b>
          <ul>
            <li>Test backend at Postman: https://www.postman.com/</li>
          </ul>
        </li>
      </ul>
    <br>
    <li &nbsp;>Part 3 - React<br/>
      <ul &nbsp;> 
        <li &nbsp;><b>Frontend project</b>
          <ul>
            <li>Create feature: modal
              <ul>
                <li>Install Radix-UI: npm i @radix-ui/react-dialog </li>
                <li>Install Radix-UI: npm i @radix-ui/react-popover </li>
                <li>Install CLSX: npm i clsx </li>
              </ul>
            </li>
            <li>Get form data</li>
            <li>Day detail popover</li>
            <li>Synchronize completed habit</li>
          </ul>
        </li>
        <li &nbsp;><b>Testing</b>
          <ul>
            <li>Test backend at Hoppscotch: https://hoppscotch.io/</li>
          </ul>
        </li>
      </ul>
    </li>
    <br>
    <li &nbsp;>Part 4 - React<br/>
      <ul &nbsp;>
        <li &nbsp;><b>Frontend project</b>
          <ul>
            <li>Use Radix UI Components 
              <ul>
                <li>Checkbox = npm install @radix-ui/react-checkbox</li>
                <li>Select = npm install @radix-ui/react-select</li>
                <li>Toggle Group = npm install @radix-ui/react-toggle-group</li>
              </ul>
            </li>
            <li>Get form data</li>
            <li>Syncing completed habits</li>
            <li>Connection with back-end
              <ul>
                <li>Send modal values to API (backend service)
                  <ul>
                    <li>Axios: npm install axios</li>
                  </ul>
                </li>
                <li>Set HTTP client</li>
                <li>Fetching API summary</li>
                <li>Using API data at Popover</li>
                <li>Register a new habit</li>
              </ul>
            </li>
          </ul>
        </li>
      </ul> 
    </li> 
    <br>
  </ol>
</div>

<div id="howtouse">
<h2>🧪 How to use</h2>
  <ol &nbsp;>
    <li &nbsp;>Set the development environment at you local computer</li>
    <li &nbsp;>Clone the repository 
      <ul>
        <li>git clone https://github.com/alexandrecpedro/ama-room</li>
      </ul>
    </li>
    <li &nbsp;>Enter the project directory: 
      <ul>
        <li>cd rocketseat-auction</li>
      </ul>
    </li>
    <li><b>Testing</b> 
      <ul>
        <li &nbsp;><u>Backend</u>
          <ol>
            <li &nbsp;>Enter the directory
              <ul>
                <li>cd backend</li>
              </ul>
            </li>
            <li &nbsp;>Install the dependencies
              <ul>
                <li>go mod tidy</li>
              </ul>
            </li>
            <li &nbsp;>Run
              <ul>
                <li>go run ./cmd/wsrs/main.go</li>
              </ul>
            </li>
          </ol>
        </li>
        <li &nbsp;><u>Frontend</u>
          <ol>
            <li &nbsp;>Enter the directory
              <ul>
                <li>cd frontend</li>
              </ul>
            </li>
            <li &nbsp;>Install the dependencies
              <ul>
                <li>npm install</li>
              </ul>
            </li>
            <li &nbsp;>Run
              <ul>
                <li>npm run dev</li>
              </ul>
            </li>
          </ol>
        </li>
      </ul>
    </li>
  </ol>
</div>
