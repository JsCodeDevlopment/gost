# 12 - Gost CLI Automation

The **Gost CLI** is the framework's automation engine. It is designed to eliminate boilerplate code and speed up the development process by scaffolding projects and generating full-stack modules following Gost's architectural standards.

---

## 📂 CLI Architecture

The CLI source code is localized in `cmd/gost` and follows a command-pattern structure using the [Cobra](https://github.com/spf13/cobra) library.

### 1. Entry Point (`main.go`)
The main entry point that simply calls `commands.Execute()`.

### 2. Root Command (`commands/root.go`)
Sets up the base `gost` command and provides the foundational CLI structure.

### 3. Project Initializer (`commands/init.go`)
The most complex command. It handles:
- **Interactive Prompts**: Uses `survey` to ask for project name and template type.
- **Directory Cloning**: Copies the current repository into a new folder, excluding development files like `.git`, `cmd`, and `test`.
- **Module Pruning**: If the user chooses a "Basic" template, it physically removes the code and configurations for modules not selected (e.g., Auth, RabbitMQ, i18n).
- **Template Patching**: Replaces the generic `gost` module name with the new project name across all `.go` and `go.mod` files.

### 4. Module Scaffolder (`commands/make_module.go`)
Generates a standard folder structure for a domain module:
- `src/modules/<name>/dto`
- `src/modules/<name>/entities`
- `src/modules/<name>/repositories`
- `src/modules/<name>/services`
- It also creates a base `<name>.module.go` file with an `InitModule` function.

### 5. CRUD Generator (`commands/make_crud.go`)
The productivity powerhouse. It performs a complete code generation cycle:
- **Entity Generation**: Creates a GORM model.
- **DTO Generation**: Creates Create/Update structs with validation tags.
- **Repository/Service/Controller**: Generates all three layers with pre-built logic for Create, FindAll, FindOne, Update, and Delete.
- **Auto-Registration**: Parses `src/app/app.module.go`, adds the necessary imports, and calls the new module's `InitModule(api)` within the `SetupApp` function.

---

## 🚀 Commands & Usage

### Initialize a New Project

**Interactive Mode:**
```bash
gost init
```

**Flag-based Mode (Scriptable):**
```bash
gost init --name my-api --template Basic --modules auth,i18n
```
*Modules available: `auth`, `messaging`, `i18n`.*

### Generate a Domain Module
```bash
gost make:module catalog
```

### Full Stack CRUD Generation
```bash
gost make:crud order
```
*Note: Pass the singular name. It will automatically pluralize the folder and endpoints (e.g., `order` -> `/api/v1/orders`).*

---

## 🔄 Internal Logic Flows

### Project Scaffolding Flow
1. **Copy**: Recursive copy of the Gost repo.
2. **Prune**: If "Basic", delete unwanted folders:
   - `!auth` -> removes `src/modules/auth`, `src/common/guards`, `src/common/security`.
   - `!messaging` -> removes `src/common/messaging`, `src/config/rabbitmq.go`.
   - `!i18n` -> removes `src/common/i18n`, `locales/`.
3. **Patch**: 
   - Update `go.mod` module name.
   - Update all imports in `.go` files.
   - Remove init calls in `app.module.go` for pruned modules.

### CRUD Generation Flow
1. **Templates**: Injected Go-string templates with placeholders.
2. **Project Detection**: Reads the current `go.mod` to ensure imports match your application name.
3. **File Creation**: Writes 6 distinct files.
4. **Registration**: 
   - Finds `import (` and injects the new module path.
   - Finds `ws.InitModule(api)` and injects the new `<module>.InitModule(api)` below it.

---

## ⚠️ Important Considerations

- **Server Restart**: After running `make:crud`, you must restart your Go server for the new routes to be registered.
- **Database Migration**: The CLI generates the Entity, but you should ensure your DB is configured to handle the new table (GORM's AutoMigrate is recommended in your module init).
- **Naming**: Always use lowercase singular names for the commands. The CLI handles capitalization and pluralization for you.
