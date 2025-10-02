# Sistema de Gestão Psicológica

## Índice

- [Visão Geral](#visão-geral)
- [Linguagem Ubíqua](#linguagem-ubíqua)
- [Arquitetura Central: O Conceito de Workspace](#arquitetura-central-o-conceito-de-workspace)
- [Casos de Uso](#casos-de-uso)
  - [Fluxo do Psicólogo Autônomo (MVP)](#fluxo-do-psicólogo-autônomo-mvp)
  - [Fluxo da Clínica (Expansão Futura)](#fluxo-da-clínica-expansão-futura)
- [Endpoints da API](#endpoints-da-api)
- [Modelo de Dados](#modelo-de-dados)
- [Architecture Decision Records (ADRs)](#architecture-decision-records-adrs)

---

## Visão Geral

O objetivo deste projeto é desenvolver um sistema de gestão focado em **psicólogos autônomos**, facilitando o controle de pacientes, agenda, prontuários e finanças.

A arquitetura foi projetada para ser simples no presente, mas escalável para o futuro. O núcleo do design é o conceito de **Workspace**, que funciona como um contêiner de dados. Para o psicólogo autônomo (MVP), isso é totalmente transparente, funcionando como seu consultório particular. No entanto, essa mesma estrutura permite, opcionalmente, uma expansão futura para suportar clínicas com múltiplos profissionais, sem a necessidade de reestruturações complexas no banco de dados.

### Principais Funcionalidades (MVP)

- **Gestão de Pacientes**: Cadastro completo com dados pessoais e de contato.
- **Gestão de Convênios**: Cadastro e associação de planos de saúde.
- **Agenda Inteligente**: Marcação de sessões únicas ou recorrentes.
- **Prontuário Eletrônico Seguro**: Registro criptografado da evolução dos pacientes.
- **Controle Financeiro Simplificado**: Gestão de pagamentos de sessões e relatórios.
- **Segurança e Privacidade**: Foco em conformidade com a LGPD e o sigilo profissional.

---

## Linguagem Ubíqua

Esta seção define os termos essenciais do sistema para garantir que todos (desenvolvedores, stakeholders e a própria documentação) falem a mesma língua.

| Termo | Definição | Contexto Técnico / Exemplo |
|-------|-----------|------------------|
| **Conta (Account)** | Representa uma pessoa que usa o sistema. Contém credenciais de login e dados pessoais. | Tabela `account`. Um psicólogo, um secretário ou um administrador têm, cada um, uma `Account`. |
| **Workspace** | A entidade central de **isolamento de dados**. É a "fronteira" que separa os dados de um consultório dos de outro. | Tabela `workspace`. É a implementação do conceito de **Tenant**. |
| **Tenant (Inquilino)** | O termo técnico para uma instância de cliente em um sistema multi-inquilino. No nosso caso, cada `Workspace` é um `Tenant`. | A estratégia de Multi-Tenancy é explicada no ADR-005. |
| **Membro (Member)** | Uma `Account` que tem acesso a um `Workspace`. A relação entre eles define o que o usuário pode fazer. | Representado pela tabela `workspace_member`. Ex: "O Dr. Carlos é um **Membro** do **Workspace** da Clínica Bem-Viver". |
| **Papel (Role)** | Define o nível de permissão de um `Membro` dentro de um `Workspace`. | Campo `role` na tabela `workspace_member`. Ex: `owner`, `psychologist`, `secretary`. |
| **Perfil (Profile)** | Não é uma entidade de banco, mas um **conceito**. É a visão unificada de todos os dados do psicólogo. | Um DTO na API que agrega dados da `account`, do `workspace`, `address`, etc. para simplificar a interação com o frontend. |
| **Profissional** | Um `Membro` com o `Papel` de `psychologist`, habilitado a realizar atendimentos. | A coluna `professional_id` na tabela `appointment` é uma FK para a tabela `account`. |
| **Paciente** | Pessoa que recebe atendimento. Sempre pertence a um único `Workspace`. | Tabela `patient`. |
| **Agendamento** | O registro de uma sessão em data e horário específicos, com status e valor. | Tabela `appointment`. |
| **Prontuário** | O conjunto de todas as evoluções clínicas de um paciente. | Conceitual. É a coleção de `progress_note`s de um `patient`. |
| **Evolução** | As anotações clínicas de uma sessão específica. O conteúdo é sempre criptografado. | Tabela `progress_note`. |
| **Série Recorrente** | Um modelo para gerar múltiplos `Agendamentos` que se repetem em um padrão (ex: semanalmente). | Tabela `recurring_series`. |
| **Soft Delete** | Prática de marcar um registro como excluído em vez de removê-lo fisicamente. | Usado em quase todas as tabelas. Ver ADR-003. |

---

## Arquitetura Central: O Conceito de Workspace

O `Workspace` é a peça fundamental da arquitetura deste sistema. Compreendê-lo é essencial para entender como os dados são organizados e protegidos.

### O que é um Workspace?

Pense no Workspace como um **contêiner digital isolado** ou uma "pasta" segura. Todos os dados operacionais — pacientes, agendamentos, prontuários, convênios — existem *dentro* de um Workspace. Ele serve como uma fronteira lógica que garante que os dados do "Consultório A" jamais se misturem com os dados do "Consultório B".

Tecnicamente, ele é a nossa implementação de **multi-tenancy** (múltiplos inquilinos), onde cada psicólogo ou clínica é um inquilino do sistema, com seus dados completamente segregados.

### Por que usar Workspaces em vez de ligar tudo à Conta?

A abordagem mais simples seria ter um `account_id` em cada tabela (`patient`, `appointment`, etc.). Isso funciona bem para um único psicólogo, mas se torna um grande problema se quisermos permitir colaboração (clínicas).

O Workspace resolve isso de forma elegante:

- **Dados pertencem ao Consultório/Clínica, não a uma Pessoa:** Pacientes são da clínica, não do psicólogo que saiu de lá. O Workspace modela essa realidade.
- **Flexibilidade de Acesso:** Múltiplas `Account`s (usuários) podem ter acesso ao mesmo `Workspace` com diferentes `Papel`s (permissões), permitindo a colaboração.

### Como funciona na Prática?

A relação entre as entidades centrais é:

```
+--------------+       +---------------------+       +---------------+
|   Account    |  (é)  |  Workspace Member   | (de)  |   Workspace   |
| (O Usuário)  |------>| (Com um Papel)      |<------| (O Contêiner) |
+--------------+       +---------------------+       +---------------+
     |                                                      |
     | (Realiza)                                            | (Contém)
     v                                                      v
+--------------+                                      +-----------+
| Appointment  |                                      |  Patient  |
+--------------+                                      +-----------+
```

#### Cenário 1: Psicólogo Autônomo (MVP)

O objetivo aqui é manter a **complexidade zero** para o usuário.

1. **Cadastro**: Quando um novo psicólogo cria sua `Account`, o sistema, **automaticamente e de forma transparente**, cria um `Workspace` do tipo `SOLO_PRACTICE` para ele.
2. **Ligação**: O sistema também cria um registro na `workspace_member` que diz: "Esta `Account` é a `owner` (dona) deste `Workspace`".
3. **Uso Diário**: Quando o psicólogo faz login, o sistema sabe qual é o seu `workspace_id`. Todas as operações (criar paciente, agendar sessão) são automaticamente associadas a esse ID. Para o usuário e para o código da API, é como se ele estivesse em seu próprio sistema privado.

#### Cenário 2: Clínica (Expansão Futura)

Aqui, o poder do Workspace se torna explícito.

1. **Criação**: O dono da clínica cria uma `Account` e, em seguida, um `Workspace` do tipo `CLINIC`.
2. **Convites**: Ele pode convidar outros profissionais e secretários. Ao aceitar, novas `Account`s são criadas (se não existirem) e registros são adicionados à `workspace_member`, ligando esses novos usuários ao Workspace da clínica com `Papel`s específicos (`psychologist`, `secretary`).
3. **Acesso Controlado**: Um secretário (`role = 'secretary'`) pode ver a agenda de todos os profissionais do Workspace, mas a API o bloqueará de ler o conteúdo dos prontuários. Um psicólogo (`role = 'psychologist'`) só verá os pacientes que lhe forem atribuídos.

Essa estrutura, definida desde o início, é o que garante que o sistema pode crescer sem a necessidade de uma migração de dados dolorosa e complexa.

---

## Casos de Uso

### Fluxo do Psicólogo Autônomo (MVP)

#### 1. Gestão de Conta e Perfil

**CreateAccount**

- **Descrição**: Um novo psicólogo se cadastra.
- **Endpoint**: `POST /auth/signup`
- **Ator**: Psicólogo (não autenticado).
- **Input**: `name`, `email`, `password`, `phone`.
- **Output**: `accountId`.
- **Regras**: `email` deve ser único.
- **Lógica de Sistema (Oculta)**:
  1. Cria uma `account` com `status = 'pending_verification'`.
  2. Cria um `workspace` do tipo `SOLO_PRACTICE` com o nome do psicólogo.
  3. Cria um vínculo na `workspace_member` com `role = 'owner'`.
  4. Envia um e-mail de verificação. A conta só se torna `active` após a verificação.

**UpdateProfile**

- **Descrição**: O psicólogo atualiza suas informações de perfil e do consultório.
- **Endpoint**: `PUT /profile`
- **Ator**: Psicólogo (autenticado).
- **Input**: Objeto `profile` contendo: `name`, `phone`, `defaultSessionValue`, `practiceName`, `address`, `contacts[]`.
- **Output**: `void`.
- **Lógica de Sistema (Oculta)**: A API recebe um único objeto e distribui as atualizações para as tabelas correspondentes: `account` (dados pessoais), `workspace` (nome do consultório), `address` e `contact`.

**DeleteAccount**

- **Descrição**: O psicólogo solicita a desativação permanente da sua conta (soft delete).
- **Endpoint**: `DELETE /account`
- **Ator**: Psicólogo (autenticado).
- **Output**: `void`.
- **Regras**: Altera o `status` da `account` e do `workspace` associado para `deleted`. Cancela agendamentos futuros. A ação é irreversível pelo usuário.

#### 2. Gestão de Pacientes

*(Nota: Em todos os casos de uso do MVP, o `workspace_id` é automaticamente inferido a partir do usuário autenticado, simplificando a lógica da API.)*

**CreatePatient**

- **Descrição**: Adicionar um novo paciente.
- **Endpoint**: `POST /patients`
- **Input**: `name`, `birthDate`, `legalGuardianName` (opcional), `contacts[]`.
- **Output**: `patientId`.
- **Regras**: `legalGuardianName` é obrigatório se a idade do paciente for menor que 18 anos.

**ListPatients**, **GetPatientDetails**, **UpdatePatient**, **DeletePatient (Soft Delete)**

- **Descrição**: Operações padrão de CRUD para pacientes, sempre restritas ao `workspace_id` do psicólogo.
- **Endpoints**: `GET /patients`, `GET /patients/:patientId`, `PUT /patients/:patientId`, `DELETE /patients/:patientId`

#### 3. Gestão de Convênios

**CreateInsurance**, **DeleteInsurance (Soft Delete)**

- **Descrição**: Adicionar e remover convênios.
- **Endpoints**: `POST /insurances`, `DELETE /insurances/:insuranceId`
- **Regras**: Um convênio não pode ser excluído se houver pacientes ativos vinculados a ele. A unicidade do nome é por `workspace`.

#### 4. Gestão de Agendamentos

**ScheduleSingleSession**

- **Descrição**: Marcar uma única sessão.
- **Endpoint**: `POST /appointments`
- **Input**: `patientId`, `scheduledDateTime`, `durationMinutes`, `sessionValue`.
- **Output**: `appointmentId`.
- **Regras**: Valida conflito de horário na agenda do profissional (`professional_id`).

**ScheduleRecurringSession**

- **Descrição**: Criar uma série de sessões recorrentes (ex: toda quarta-feira às 10h).
- **Endpoint**: `POST /appointments`
- **Output**: `seriesId`.
- **Lógica**: Gera todas as instâncias de `appointment` no momento da criação (ADR-002).

**RescheduleSession**

- **Descrição**: Operações de gerenciamento de agendamentos.
- **Endpoint**: `PUT /appointments/:appointmentId/reschedule`

**CancelSession**

- **Descrição**: Operações de gerenciamento de agendamentos.
- **Endpoint**: `DELETE /appointments/:appointmentId`

**CancelRecurringSeries**

- **Descrição**: Operações de gerenciamento de agendamentos.
- **Endpoint**: `DELETE /recurring-series/:seriesId`

#### 5. Gestão de Prontuários

**CreateProgressNote**

- **Descrição**: Registrar a evolução de uma sessão concluída.
- **Endpoint**: `POST /progress-notes`
- **Input**: `appointmentId`, `sessionSummary`.
- **Output**: `progressNoteId`.
- **Regras**:
  - Requer que o `appointment` tenha `status = 'completed'`.
  - O campo `sessionSummary` **deve ser criptografado** pela aplicação antes de ser salvo (ADR-001).

**UpdateProgressNote**

- **Descrição**: Editar uma evolução já registrada.
- **Endpoint**: `PUT /progress-notes/:noteId`
- **Regras**: A edição é bloqueada pela aplicação 30 dias após a data de criação da nota (ADR-004).

**GetPatientClinicalHistory**

- **Descrição**: Visualizar o prontuário completo de um paciente.
- **Endpoint**: `GET /patients/:patientId/clinical-history`
- **Lógica**: A aplicação busca todas as `progress_note`s e **decripta o `sessionSummary`** antes de retornar os dados para a interface.

---

### Fluxo da Clínica (Expansão Futura)

#### 1. Gestão da Clínica e Membros

**CreateClinicWorkspace**, **InviteMember**, **ManageMember**

- **Descrição**: Funcionalidades administrativas para o dono da clínica criar o workspace, convidar profissionais e gerenciar seus papéis.
- **Endpoints**: `POST /workspaces`, `POST /workspaces/:workspaceId/members`, `PUT /workspaces/:workspaceId/members/:memberId`

**SwitchWorkspaceContext**

- **Descrição**: Um usuário que pertence a múltiplos workspaces (seu consultório particular e uma clínica) pode alternar entre eles na interface.
- **Endpoint**: Não é um endpoint. A lógica é implementada no cliente, que passa a enviar um header (ex: `X-Workspace-ID`) nas requisições.
- **Lógica**: A aplicação (frontend) envia um header (ex: `X-Workspace-ID`) nas requisições. A API (backend) usa esse ID para filtrar todas as consultas, aplicando as permissões do papel (`role`) do usuário *naquele* contexto.

#### 2. Operações Diárias na Clínica

**CreateClinicPatient**, **ScheduleSessionForProfessional**, **ManageSessionPayment**

- **Descrição**: Operações similares ao MVP, mas com a complexidade adicional de múltiplos profissionais e papéis (ex: um secretário agendando para um psicólogo).
- **Endpoints**: `POST /patients`, `POST /appointments`, `PUT /appointments/:appointmentId/payment`

#### 3. Controle de Acesso e Visibilidade

**AccessPatientClinicalHistory**

- **Descrição**: Acessar o prontuário de um paciente da clínica.
- **Endpoint**: `GET /patients/:patientId/clinical-history`
- **REGRAS DE NEGÓCIO CRÍTICAS**:
  - A API deve garantir que apenas o(s) psicólogo(s) associado(s) ao tratamento do paciente possam visualizar o conteúdo do prontuário.
  - Papéis administrativos (`admin`, `secretary`) **NUNCA** devem ter acesso ao conteúdo descriptografado das evoluções (`session_summary`). Essa é uma restrição de segurança e privacidade inegociável.

---

## Endpoints da API

A API para o MVP é projetada para ser simples e centrada no conceito de "perfil" do psicólogo.

```http
# --- Autenticação e Conta ---
POST   /auth/signup                   # Cria conta, workspace e envia email de verificação
POST   /auth/login
POST   /auth/verify-email
POST   /auth/resend-verification
POST   /auth/forgot-password
POST   /auth/reset-password

# --- Perfil e Configurações (opera no contexto do usuário/workspace atual) ---
GET    /profile                       # Agrega e retorna dados de account, workspace, etc.
PUT    /profile                       # Atualiza o perfil completo do psicólogo
DELETE /account                       # Inicia o processo de soft delete da conta

# --- Entidades do Domínio (sempre filtradas pelo workspace_id da sessão) ---
GET    /patients
POST   /patients
GET    /patients/:patientId
PUT    /patients/:patientId
DELETE /patients/:patientId

GET    /insurances
POST   /insurances
DELETE /insurances/:insuranceId

# ... e assim por diante para Appointments, Progress Notes, etc.

# --- Relatórios ---
GET    /reports/schedule?startDate=...&endDate=...
GET    /reports/financial?startDate=...&endDate=...
```

---

## Modelo de Dados

```sql
-- Extensão para gerar UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS psychological_management;
COMMENT ON SCHEMA psychological_management IS 'Schema para todas as tabelas do sistema de gestão psicológica.';

-- ===================================================================
-- 1. ENTIDADES CENTRAIS: Workspace e Account
-- ===================================================================

CREATE TABLE psychological_management.workspace (
    workspace_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_type VARCHAR(20) NOT NULL CHECK (workspace_type IN ('SOLO_PRACTICE', 'CLINIC')),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE psychological_management.workspace IS 'Entidade de tenant que isola todos os dados. A raiz da estratégia de multi-tenancy.';
COMMENT ON COLUMN psychological_management.workspace.name IS 'Para SOLO_PRACTICE, é o "Nome do Consultório". Para CLINIC, o nome da clínica.';

CREATE TABLE psychological_management.account (
    account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    default_session_value DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'pending_verification' CHECK (status IN ('pending_verification', 'active', 'inactive', 'locked', 'suspended', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE psychological_management.account IS 'Representa um usuário individual e suas credenciais de login.';
COMMENT ON COLUMN psychological_management.account.status IS 'pending_verification: aguardando confirmação de email. active: normal. inactive: desativado pelo usuário. locked: bloqueado por segurança (ex: falha de login). suspended: bloqueado por um admin. deleted: marcado para exclusão.';

CREATE TABLE psychological_management.workspace_member (
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES psychological_management.account(account_id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'admin', 'psychologist', 'secretary', 'financial')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, account_id)
);
COMMENT ON TABLE psychological_management.workspace_member IS 'Tabela de junção que define o papel de uma Account em um Workspace.';

-- ===================================================================
-- 2. TABELAS POLIMÓRFICAS: Address e Contact (reutilizáveis)
-- ===================================================================

CREATE TABLE psychological_management.address (
    address_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    owner_type VARCHAR(20) NOT NULL CHECK (owner_type IN ('workspace', 'patient')),
    street VARCHAR(255),
    -- ... outros campos de endereço
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE psychological_management.address IS 'Endereços polimórficos. owner_id/owner_type apontam para a entidade dona do endereço (ver ADR-007).';

CREATE TABLE psychological_management.contact (
    contact_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    owner_type VARCHAR(20) NOT NULL CHECK (owner_type IN ('workspace', 'patient')),
    contact_type VARCHAR(20) NOT NULL CHECK (contact_type IN ('phone_mobile', 'email')),
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE psychological_management.contact IS 'Contatos polimórficos (telefones, emails) para diferentes entidades.';
CREATE INDEX idx_contact_owner ON psychological_management.contact(owner_id, owner_type);

-- ===================================================================
-- 3. ENTIDADES DO DOMÍNIO (sempre ligadas a um Workspace)
-- ===================================================================

CREATE TABLE psychological_management.insurance (
    insurance_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id),
    name VARCHAR(255) NOT NULL,
    session_value DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (workspace_id, name)
);
COMMENT ON TABLE psychological_management.insurance IS 'Convênios/Planos de Saúde cadastrados por um Workspace.';

CREATE TABLE psychological_management.patient (
    patient_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id),
    name VARCHAR(255) NOT NULL,
    birth_date DATE,
    legal_guardian_name VARCHAR(255),
    insurance_id UUID REFERENCES psychological_management.insurance(insurance_id) ON DELETE SET NULL,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'discharged', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE psychological_management.patient IS 'Pacientes do consultório ou clínica.';

CREATE TABLE psychological_management.appointment (
    appointment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id),
    professional_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    patient_id UUID NOT NULL REFERENCES psychological_management.patient(patient_id),
    series_id UUID, -- FK para recurring_series, se aplicável
    scheduled_datetime TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER DEFAULT 50,
    session_value DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'no_show', 'cancelled')),
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE psychological_management.appointment IS 'Registro de uma sessão agendada.';
COMMENT ON COLUMN psychological_management.appointment.professional_id IS 'Indica qual profissional (Account) está conduzindo a sessão.';
CREATE INDEX idx_appointment_professional_time ON psychological_management.appointment(professional_id, scheduled_datetime);
CREATE INDEX idx_appointment_workspace_time ON psychological_management.appointment(workspace_id, scheduled_datetime);

CREATE TABLE psychological_management.progress_note (
    progress_note_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    appointment_id UUID NOT NULL REFERENCES psychological_management.appointment(appointment_id) UNIQUE,
    session_summary TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE psychological_management.progress_note IS 'Anotação da evolução clínica de uma sessão.';
COMMENT ON COLUMN psychological_management.progress_note.session_summary IS 'Conteúdo SENSÍVEL. Deve ser criptografado em nível de aplicação (ver ADR-001).';
```

---

## Architecture Decision Records (ADRs)

Registros das decisões de arquitetura mais importantes e suas consequências.

### ADR 001: Armazenar Dados Sensíveis com Criptografia

- **Decisão**: Implementar criptografia em nível de aplicação (ex: AES-256) para campos sensíveis como `session_summary`. As chaves de criptografia serão gerenciadas por um serviço de segredos (ex: AWS KMS, HashiCorp Vault), separado do banco de dados e do código-fonte.
- **Consequências**: ✅ Segurança robusta e conformidade com a LGPD. ❌ Aumenta a complexidade da aplicação (gerenciamento de chaves) e a latência de leitura/escrita desses dados.

### ADR 002: Geração de Sessões Recorrentes

- **Decisão**: No momento da criação de uma série recorrente, gerar e salvar todas as ocorrências futuras como registros individuais na tabela `appointment`. Uma tabela `recurring_series` (não detalhada aqui para simplicidade) pode ser usada para agrupar e gerenciar a série como um todo.
- **Consequências**: ✅ Simplifica drasticamente a consulta da agenda (um simples `SELECT` em um período). Permite alterações e cancelamentos individuais em sessões de uma série. ❌ Aumenta o uso de armazenamento. É necessário impor um limite razoável de ocorrências (ex: 1 ano) para evitar sobrecarga.

### ADR 003: Soft Delete vs Hard Delete

- **Decisão**: Implementar soft delete usando um campo `status` (ex: `deleted`) e/ou um timestamp `deleted_at`. A lógica de negócio na aplicação será responsável por filtrar esses registros em todas as consultas.
- **Consequências**: ✅ Mantém a rastreabilidade e integridade referencial, essencial para auditoria em sistemas de saúde. ❌ Exige disciplina na implementação: toda query deve incluir a cláusula de filtro (ex: `WHERE status != 'deleted'`), o que pode ser propenso a erros.

### ADR 004: Controle de Edição de Evoluções

- **Decisão**: A regra de negócio que impede a edição de `progress_note` após 30 dias será implementada na camada de aplicação/serviço, não no banco de dados.
- **Consequências**: ✅ Garante conformidade com as normas profissionais (CFP). Manter a lógica na aplicação facilita a criação de testes unitários e a manutenção da regra, caso ela mude no futuro.

### ADR 005: Isolamento de Dados via Workspace (Multi-Tenancy)

- **Decisão**: Implementar uma estratégia de Multi-Tenancy lógica. A entidade `workspace` é a raiz do tenant. Todas as tabelas principais (patient, appointment, etc.) terão uma Foreign Key `workspace_id`. A camada de acesso a dados da aplicação (ex: um repositório base ou middleware) será responsável por adicionar a cláusula `WHERE workspace_id = :current_workspace_id` a todas as consultas automaticamente.
- **Consequências**: ✅ Solução simples e eficaz para isolar dados, com baixo custo de implementação e manutenção. ❌ Depende criticamente da implementação correta na aplicação para garantir a segurança. Testes rigorosos são necessários para prevenir vazamento de dados entre tenants.

### ADR 006: Modelagem de Workspace para Múltiplos Tipos

- **Decisão**: Usar uma única tabela `workspace` com uma coluna `workspace_type` para diferenciar consultórios autônomos (`SOLO_PRACTICE`) de clínicas (`CLINIC`). Dados específicos de clínicas (como CNPJ) poderiam ficar em uma tabela separada (`clinic_profile`) seguindo o padrão "Class Table Inheritance".
- **Consequências**: ✅ Esquema limpo e normalizado que evita colunas nulas na tabela principal. Facilita a expansão para novos tipos de workspace no futuro.

### ADR 007: Modelagem Polimórfica para Endereços e Contatos

- **Decisão**: Usar tabelas `address` e `contact` com associação polimórfica (colunas `owner_id` e `owner_type`) para evitar a duplicação de estruturas.
- **Consequências**: ✅ Reutilização de código e estrutura (DRY). ❌ Aumenta a complexidade das queries (requer `JOIN`s com condições). O SGBD não pode garantir a integridade referencial com uma FK nativa, então a lógica de consistência fica a cargo da aplicação.
