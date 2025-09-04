# Sistema de Gestão Psicológica

## Índice
- [Visão Geral](#visão-geral)
- [Linguagem Ubíqua](#linguagem-ubíqua)
- [Casos de Uso](#casos-de-uso)
  - [Fluxo do Psicólogo Autônomo (MVP)](#fluxo-do-psicólogo-autônomo-mvp)
  - [Fluxo da Clínica (Expansão Futura)](#fluxo-da-clínica-expansão-futura)
- [Endpoints da API](#endpoints-da-api)
- [Modelo de Dados](#modelo-de-dados)
- [Architecture Decision Records](#architecture-decision-records)

---

## Visão Geral

O objetivo deste projeto é desenvolver um sistema de gestão focado em **psicólogos autônomos**, facilitando o controle de pacientes, agenda, prontuários e aspectos financeiros.

A arquitetura foi projetada para ser escalável, permitindo uma futura e opcional expansão para suportar clínicas com múltiplos profissionais, sem a necessidade de reestruturações complexas na base de dados. O foco inicial e principal, no entanto, é a simplicidade e utilidade para o profissional autônomo.

### Principais Funcionalidades

- **Gestão de Pacientes**: Cadastro completo com dados pessoais e de contato.
- **Sistema de Convênios**: Cadastro e gestão de planos de saúde.
- **Agendamentos Inteligentes**: Sessões únicas ou recorrentes.
- **Prontuário Eletrônico**: Registro seguro da evolução dos pacientes.
- **Controle Financeiro**: Gestão de pagamentos e relatórios.
- **Segurança e Privacidade**: Conformidade com LGPD e sigilo profissional.

---

## Linguagem Ubíqua

| Termo | Definição |
|-------|-----------|
| **Conta (Account)** | Representa o psicólogo (profissional licenciado) que utiliza o sistema. Contém seus dados pessoais e configurações. |
| **Perfil (Profile)** | A visão unificada dos dados do psicólogo, incluindo informações pessoais, configurações e detalhes do consultório. |
| **Workspace** | (Termo Interno) A entidade que isola os dados de um psicólogo ou clínica. Totalmente transparente para o usuário no MVP. |
| **Unidade (Location)** | (Futuro) Uma filial ou endereço físico de uma Clínica. |
| **Paciente** | Pessoa que recebe atendimento psicológico. |
| **Sessão** | Encontro terapêutico entre psicólogo e paciente. |
| **Agendamento** | Marcação de uma sessão em data e horário específicos. |
| **Prontuário** | Registro clínico digital contendo histórico e evolução do paciente. |
| **Evolução** | Anotações clínicas sobre o desenvolvimento do paciente. |
| **Soft Delete** | Marcação de registro como excluído sem remoção física dos dados. |

---

## Casos de Uso

### Fluxo do Psicólogo Autônomo (MVP)

Esta seção detalha o escopo inicial e principal do projeto. O conceito de "Workspace" é um detalhe de implementação e deve ser totalmente transparente para o usuário.

#### **1. Gestão de Conta e Perfil**

**CreateAccount**
- **Descrição**: Um novo psicólogo se cadastra para usar o sistema.
- **Ator**: Psicólogo (não autenticado).
- **Input**: `name`, `email`, `password`, `phone`.
- **Output**: `accountId`.
- **Regras de Negócio**: `email` deve ser único no sistema.
- **Lógica de Sistema (Oculta)**: Cria um `account`, um `workspace` do tipo `SOLO_PRACTICE` e um vínculo `workspace_member` com `role = 'owner'`.

**UpdateProfile**
- **Descrição**: O psicólogo atualiza todas as informações editáveis do seu perfil e consultório.
- **Ator**: Psicólogo (autenticado).
- **Input**: Objeto `profile` contendo: `name`, `phone`, `defaultSessionValue`, `practiceName`, `address`, `contacts[]`.
- **Output**: `void`.
- **Lógica de Sistema (Oculta)**: Atualiza as tabelas `account` (dados pessoais), `workspace` (nome do consultório), `address` e `contact` (dados do consultório).

**DeactivateAccount**
- **Descrição**: O psicólogo desativa sua conta (soft delete).
- **Ator**: Psicólogo (autenticado).
- **Output**: `void`.
- **Regras**: Altera o status da `account` e do `workspace` associado para `inactive` e cancela agendamentos futuros.

#### **2. Gestão de Pacientes**

**CreatePatient**
- **Descrição**: Adicionar um novo paciente ao consultório.
- **Ator**: Psicólogo (autenticado).
- **Input**: `name`, `birthDate`, `legalGuardianName` (opcional), `contacts[]`.
- **Output**: `patientId`.
- **Regras**: `legalGuardianName` é obrigatório se o paciente for menor de 18 anos.

**ListPatients**
- **Descrição**: Visualizar a lista de pacientes do consultório.
- **Ator**: Psicólogo (autenticado).
- **Input**: `status` (filtro: 'active', 'inactive'), `page`.
- **Output**: Lista paginada de pacientes.

**GetPatientDetails**
- **Descrição**: Acessar o perfil completo de um paciente.
- **Ator**: Psicólogo (autenticado).
- **Input**: `patientId`.
- **Output**: Objeto completo do paciente.

**UpdatePatient**
- **Descrição**: Editar as informações de um paciente.
- **Ator**: Psicólogo (autenticado).
- **Input**: `patientId`, `name`, `contacts[]`, `insuranceId`, `status`, `notes`.
- **Output**: `void`.

**DeletePatient**
- **Descrição**: Remover um paciente (soft delete).
- **Ator**: Psicólogo (autenticado).
- **Input**: `patientId`.
- **Output**: `void`.

#### **3. Gestão de Convênios**

**CreateInsurance**
- **Descrição**: Adicionar um novo plano de saúde.
- **Ator**: Psicólogo (autenticado).
- **Input**: `name`, `sessionValue`.
- **Output**: `insuranceId`.
- **Regras**: O nome do convênio deve ser único para o psicólogo.

**DeleteInsurance**
- **Descrição**: Remover um convênio (soft delete).
- **Ator**: Psicólogo (autenticado).
- **Input**: `insuranceId`.
- **Output**: `void`.
- **Regras**: Não pode ser excluído se houver pacientes ativos vinculados.

#### **4. Gestão de Agendamentos**

**ScheduleSingleSession**
- **Descrição**: Marcar uma única sessão para um paciente.
- **Ator**: Psicólogo (autenticado).
- **Input**: `patientId`, `scheduledDateTime`, `durationMinutes`, `sessionValue`.
- **Output**: `appointmentId`.
- **Regras**: Validar conflito de horário na agenda.

**ScheduleRecurringSession**
- **Descrição**: Criar uma série de sessões recorrentes.
- **Ator**: Psicólogo (autenticado).
- **Input**: `patientId`, `recurrenceType`, `startDate`, `endDate` (ou `maxOccurrences`).
- **Output**: `seriesId`.

**RescheduleSession**
- **Descrição**: Alterar a data/hora de um agendamento.
- **Ator**: Psicólogo (autenticado).
- **Input**: `appointmentId`, `newScheduledDateTime`.
- **Output**: `void`.
- **Regras**: Validar que a nova data/hora não gera conflito.

**CancelSession**
- **Descrição**: Cancelar um único agendamento.
- **Ator**: Psicólogo (autenticado).
- **Input**: `appointmentId`, `reason` (opcional).
- **Output**: `void`.
- **Regras**: O `status` do agendamento é alterado para `cancelled` e `payment_status` também para `cancelled`.

**CancelRecurringSeries**
- **Descrição**: Cancelar todas as sessões futuras de uma série.
- **Ator**: Psicólogo (autenticado).
- **Input**: `seriesId`.
- **Output**: `cancelledSessionsCount`.
- **Regras**: Apenas agendamentos com `status = 'scheduled'` pertencentes à série são cancelados.

#### **5. Gestão de Prontuários**

**CreateProgressNote**
- **Descrição**: Registrar o resumo e as anotações de uma sessão realizada.
- **Ator**: Psicólogo (autenticado).
- **Input**: `appointmentId`, `sessionSummary`.
- **Output**: `progressNoteId`.
- **Regras**:
  - Só é possível criar uma evolução para um agendamento com `status = 'completed'`.
  - O campo `sessionSummary` deve ser criptografado antes de ser salvo (ADR-001).

**UpdateProgressNote**
- **Descrição**: Editar uma anotação de evolução já criada.
- **Ator**: Psicólogo (autenticado).
- **Input**: `progressNoteId`, `sessionSummary`.
- **Output**: `void`.
- **Regras**: A edição é permitida apenas por 30 dias após a data de criação da nota (ADR-004).

**GetPatientClinicalHistory**
- **Descrição**: Visualizar todas as evoluções de um paciente em ordem cronológica.
- **Ator**: Psicólogo (autenticado).
- **Input**: `patientId`.
- **Output**: Lista de `progress_note`s (com o `sessionSummary` decriptado).

#### **6. Relatórios**

**GetSchedule**
- **Descrição**: Visualizar os agendamentos em um período.
- **Ator**: Psicólogo (autenticado).
- **Input**: `startDate`, `endDate`.
- **Output**: Lista de `appointment`s no período, contendo dados do paciente.

**GetFinancialReport**
- **Descrição**: Resumo financeiro do consultório em um período.
- **Ator**: Psicólogo (autenticado).
- **Input**: `startDate`, `endDate`.
- **Output**: Objeto com `totalRevenue` (soma de sessões 'completed'), `paidSessions`, `pendingSessions`.

---

### Fluxo da Clínica (Expansão Futura)

Esta seção detalha os casos de uso que seriam implementados se o suporte a clínicas for adicionado. Eles dependem da arquitetura já definida (Workspaces, Membros, Papéis).

#### **1. Gestão da Clínica e Membros**

**CreateClinicWorkspace**
- **Descrição**: Um usuário (dono) cadastra uma nova clínica no sistema.
- **Ator**: Usuário (`account`) autenticado.
- **Input**: `clinicName`, `cnpj`, `address`, `contacts[]`.
- **Output**: `workspaceId`.
- **Lógica de Sistema**: Cria um `workspace` com `workspace_type = 'CLINIC'`, um `clinic_profile` com os dados fiscais e associa o usuário criador como `owner` na tabela `workspace_member`.

**InviteMember**
- **Descrição**: Um administrador convida um novo usuário (psicólogo, secretário) para a clínica.
- **Ator**: Dono (`owner`) ou Administrador (`admin`) da clínica.
- **Input**: `email`, `role` (`psychologist`, `secretary`, `financial`, `admin`).
- **Output**: `void`.
- **Regras**: O sistema envia um convite para o email. Se o usuário não existir, ele é instruído a criar uma `account`. Ao aceitar, um novo registro é criado em `workspace_member`.

**ManageMember**
- **Descrição**: Alterar o papel ou remover um membro da clínica.
- **Ator**: Dono (`owner`) ou Administrador (`admin`).
- **Input**: `memberAccountId`, `newRole` (opcional), `action` ('update' ou 'remove').
- **Output**: `void`.

**SwitchWorkspaceContext**
- **Descrição**: Um usuário que pertence a múltiplos workspaces (seu consultório particular e uma clínica) pode alternar entre eles na interface.
- **Ator**: Usuário (autenticado).
- **Lógica**: A aplicação passa a usar o `workspaceId` selecionado para todas as operações subsequentes, aplicando as permissões do `role` daquele contexto.

#### **2. Operações Diárias na Clínica**

**CreateClinicPatient**
- **Descrição**: Um secretário ou psicólogo cadastra um paciente para a clínica.
- **Ator**: `secretary`, `admin`, `psychologist`.
- **Input**: `name`, `birthDate`, `contacts[]`, `primaryProfessionalId` (opcional).
- **Output**: `patientId`.
- **Regras**: O paciente é criado dentro do `workspaceId` da clínica.

**ScheduleSessionForProfessional**
- **Descrição**: Um secretário agenda uma sessão para um dos psicólogos da clínica.
- **Ator**: `secretary`, `admin`.
- **Input**: `patientId`, **`professionalId`**, `scheduledDateTime`, `durationMinutes`.
- **Output**: `appointmentId`.
- **Regras**: O sistema valida conflitos na agenda do `professionalId` especificado.

**ManageSessionPayment**
- **Descrição**: O setor financeiro ou a secretaria atualiza o status de pagamento de uma sessão.
- **Ator**: `secretary`, `financial`, `admin`.
- **Input**: `appointmentId`, `paymentStatus` (`paid`, `pending`).
- **Output**: `void`.

#### **3. Controle de Acesso e Visibilidade**

**GetClinicSchedule**
- **Descrição**: Visualizar a agenda da clínica, com filtros por profissional.
- **Ator**: `owner`, `admin`, `secretary`.
- **Input**: `startDate`, `endDate`, `professionalId` (opcional).
- **Output**: Lista de agendamentos.

**GetMyScheduleInClinic**
- **Descrição**: Um psicólogo visualiza apenas a sua própria agenda dentro do contexto da clínica.
- **Ator**: `psychologist`.
- **Input**: `startDate`, `endDate`.
- **Output**: Lista de seus `appointment`s.

**AccessPatientClinicalHistory**
- **Descrição**: Acessar o prontuário de um paciente da clínica.
- **Ator**: `psychologist`.
- **Input**: `patientId`.
- **Output**: Histórico clínico do paciente.
- **Regras de Negócio CRÍTICAS**:
  - Apenas o(s) psicólogo(s) diretamente associado(s) ao tratamento do paciente podem visualizar o histórico.
  - Papéis administrativos (`admin`, `secretary`, `financial`) **NUNCA** devem ter acesso ao conteúdo das evoluções (`session_summary`). A API deve impor essa restrição rigorosamente.

#### **4. Relatórios da Clínica**

**GetClinicFinancialReport**
- **Descrição**: Gerar um relatório financeiro consolidado da clínica.
- **Ator**: `owner`, `admin`, `financial`.
- **Input**: `startDate`, `endDate`, `professionalId` (opcional).
- **Output**: Relatório com receita total, sessões pagas e pendentes, podendo ser quebrado por profissional.

---

## Endpoints da API

A API para o MVP é centrada no psicólogo. O conceito de "perfil" unifica as configurações.

```http
# Conta e Perfil
POST   /signup
GET    /profile             # Agrega e retorna todos os dados do perfil do psicólogo
PUT    /profile             # Atualiza o perfil completo do psicólogo
DELETE /account             # Desativa a conta

# Pacientes
POST   /patients
GET    /patients
PUT    /patients/:patientId
DELETE /patients/:patientId

# Convênios
POST   /insurances
GET    /insurances
DELETE /insurances/:insuranceId

# Agendamentos
POST   /appointments
PUT    /appointments/:appointmentId/reschedule
DELETE /appointments/:appointmentId
DELETE /recurring-series/:seriesId

# Prontuários
POST   /progress-notes
PUT    /progress-notes/:noteId

# Relatórios
GET    /schedule
GET    /reports/financial
```

---

## Modelo de Dados

```sql
-- Extensão para gerar UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS psychological_management;

-- ===================================================================
-- 1. ENTIDADES CENTRAIS: Workspace e Account
-- ===================================================================

CREATE TABLE psychological_management.workspace (
    workspace_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_type VARCHAR(20) NOT NULL CHECK (workspace_type IN ('SOLO_PRACTICE', 'CLINIC')),
    name VARCHAR(255) NOT NULL, -- No MVP, este é o "practiceName"
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE psychological_management.workspace IS 'Entidade interna de tenant. Para o autônomo, representa seu consultório.';

CREATE TABLE psychological_management.clinic_profile (
    workspace_id UUID PRIMARY KEY REFERENCES psychological_management.workspace(workspace_id) ON DELETE CASCADE,
    cnpj VARCHAR(18) UNIQUE,
    legal_representative_name VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE psychological_management.clinic_profile IS 'Dados exclusivos de workspaces do tipo CLINIC (Padrão Class Table Inheritance).';

CREATE TABLE psychological_management.account (
    account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    default_session_value DECIMAL(10,2), -- Associado diretamente ao profissional.
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE psychological_management.account IS 'Conta do usuário e local de suas configurações pessoais e financeiras.';

CREATE TABLE psychological_management.workspace_member (
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id),
    account_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'admin', 'psychologist', 'secretary', 'financial')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, account_id)
);
COMMENT ON TABLE psychological_management.workspace_member IS 'Define o papel de um usuário (account) em um workspace.';

-- ===================================================================
-- 2. TABELAS POLIMÓRFICAS: Address e Contact
-- ===================================================================

CREATE TABLE psychological_management.address (
    address_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    owner_type VARCHAR(20) NOT NULL CHECK (owner_type IN ('workspace', 'account', 'patient')),
    label VARCHAR(100) NOT NULL DEFAULT 'Principal',
    street VARCHAR(255) NOT NULL,
    number VARCHAR(20),
    complement VARCHAR(100),
    neighborhood VARCHAR(100),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(2) NOT NULL,
    zip_code VARCHAR(9) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE psychological_management.address IS 'Endereços polimórficos. O "label" diferencia múltiplos endereços para o mesmo dono (ex: Unidade Paulista).';

CREATE TABLE psychological_management.contact (
    contact_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    owner_type VARCHAR(20) NOT NULL CHECK (owner_type IN ('workspace', 'account', 'patient')),
    contact_type VARCHAR(20) NOT NULL CHECK (contact_type IN ('phone_mobile', 'phone_landline', 'email')),
    value VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE psychological_management.contact IS 'Contatos polimórficos (telefones, emails).';
CREATE INDEX idx_contact_owner ON psychological_management.contact(owner_id, owner_type);

-- ===================================================================
-- 3. ENTIDADES DO DOMÍNIO PRINCIPAL (Sempre ligadas a um Workspace)
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

CREATE TABLE psychological_management.recurring_series (
    series_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id),
    professional_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    patient_id UUID NOT NULL REFERENCES psychological_management.patient(patient_id),
    recurrence_type VARCHAR(20) NOT NULL CHECK (recurrence_type IN ('daily', 'weekly', 'monthly')),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    max_occurrences INTEGER,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'cancelled', 'completed')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE psychological_management.appointment (
    appointment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES psychological_management.workspace(workspace_id),
    professional_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    patient_id UUID NOT NULL REFERENCES psychological_management.patient(patient_id),
    series_id UUID REFERENCES psychological_management.recurring_series(series_id) ON DELETE SET NULL,
    scheduled_datetime TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER DEFAULT 50,
    session_value DECIMAL(10,2) NOT NULL,
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('private', 'insurance')),
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'no_show', 'cancelled')),
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
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

```
---
## Architecture Decision Records

### ADR 001: Armazenar Dados Sensíveis com Criptografia
- **Decisão**: Implementar criptografia em nível de aplicação (ex: AES-256) para campos sensíveis como `session_summary` na tabela `progress_note`. As chaves serão gerenciadas por um serviço seguro, separado do banco de dados.
- **Consequências**: ✅ Maior segurança e conformidade legal. ❌ Aumenta a complexidade da aplicação e a latência de leitura/escrita.

### ADR 002: Geração de Sessões Recorrentes
- **Decisão**: Gerar todos os agendamentos de uma série recorrente no momento da sua criação e armazená-los individualmente na tabela `appointment`. A tabela `recurring_series` servirá como um registro mestre para agrupar essas sessões.
- **Consequências**: ✅ Simplifica a consulta da agenda e permite modificações individuais em sessões de uma série. ❌ Aumenta o uso de armazenamento. Requer um limite para evitar sobrecarga.

### ADR 003: Soft Delete vs Hard Delete
- **Decisão**: Implementar soft delete usando um campo `status` e/ou um timestamp `deleted_at` em entidades principais. A lógica de negócio na aplicação será responsável por filtrar registros marcados como excluídos.
- **Consequências**: ✅ Mantém a rastreabilidade e permite a recuperação de dados para fins de auditoria. ❌ As queries da aplicação devem sempre incluir o filtro, aumentando o risco de erro humano.

### ADR 004: Controle de Edição de Evoluções
- **Decisão**: A lógica para controlar a edição de `progress_note` será implementada na camada de aplicação, bloqueando a edição após 30 dias da criação, conforme normas do CFP.
- **Consequências**: ✅ Garante conformidade com as normas profissionais e mantém a lógica de negócio na aplicação, facilitando testes.

### ADR 005: Isolamento de Dados via Workspace (Multi-Tenancy)
- **Decisão**: Implementar uma estratégia de Multi-Tenancy lógica. A entidade `workspace` é a raiz do tenant. Todas as tabelas de dados principais (patient, appointment, etc.) conterão uma coluna `workspace_id`. A camada de acesso a dados da aplicação será responsável por adicionar a cláusula `WHERE workspace_id = :current_workspace_id` a todas as consultas.
- **Consequências**: ✅ Solução simples e eficaz para isolar dados de diferentes consultórios/clínicas. ❌ Depende criticamente da disciplina de implementação na aplicação. Requer testes rigorosos para prevenir vazamento de dados.

### ADR 006: Modelagem de Workspace com Herança
- **Decisão**: Modelar consultórios e clínicas usando o padrão "Class Table Inheritance". A tabela base `workspace` contém os dados comuns, enquanto a tabela `clinic_profile` contém os dados específicos de uma clínica. A coluna `workspace.workspace_type` diferencia os tipos.
- **Consequências**: ✅ Esquema limpo e normalizado. Evita FKs nulas e complexidade nas tabelas de dados (`patient`, `appointment`). Facilita a expansão para novos tipos de workspace no futuro.

### ADR 007: Modelagem Polimórfica para Endereços e Contatos
- **Decisão**: Usar tabelas `address` e `contact` com associação polimórfica (colunas `owner_id` e `owner_type`) para evitar a duplicação de estruturas de endereço/contato para cada entidade (workspace, paciente, etc.).
- **Consequências**: ✅ Reutilização de estrutura (DRY). ❌ Aumenta a complexidade das queries para buscar esses dados, pois o SGBD não pode garantir a integridade referencial da FK polimórfica. A lógica fica na aplicação.

### ADR 008: Preparação para Múltiplas Unidades (Locations)
- **Decisão**: Para o MVP, não será criada uma tabela `location`. Em vez disso, a tabela `address` foi preparada para suportar múltiplos endereços para um mesmo `workspace` (removendo a `UNIQUE KEY` de `owner_id` e adicionando um campo `label`).
- **Consequências**: ✅ Mantém a simplicidade do modelo inicial. ❌ Não permite associar lógicas (profissionais, agendas) a uma unidade específica. No futuro, quando a funcionalidade for necessária, será criada a tabela `location` e os dados da tabela `address` serão migrados para o novo modelo. Isso representa um débito técnico consciente e gerenciado.- **Sistema de Convênios**: Cadastro e gestão de planos de saúde.
- **Agendamentos Inteligentes**: Sessões únicas ou recorrentes.
- **Prontuário Eletrônico**: Registro de evolução dos pacientes.
- **Controle Financeiro**: Gestão de pagamentos e relatórios.
- **Segurança e Privacidade**: Conformidade com LGPD e sigilo profissional.

---

## Linguagem Ubíqua

| Termo | Definição |
|-------|-----------|
| **Conta** | Representa o psicólogo (profissional licenciado) que utiliza o sistema. |
| **Paciente** | Pessoa que recebe atendimento psicológico. |
| **Sessão** | Encontro terapêutico entre psicólogo e paciente. |
| **Agendamento** | Marcação de uma sessão em data e horário específicos. |
| **Prontuário** | Registro clínico digital contendo histórico e evolução do paciente. |
| **Evolução** | Anotações clínicas sobre o desenvolvimento do paciente. |
| **Status do Paciente** | Situação atual: `ativo`, `inativo`, `alta`. |
| **Soft Delete** | Marcação de registro como excluído sem remoção física dos dados. |
| **Série Recorrente** | Conjunto de sessões agendadas seguindo um padrão. |

---

## Casos de Uso

### Gestão de Conta

#### CreateAccount (Criar Conta)
Criar uma conta para o psicólogo utilizar o sistema.
**Input**: `name`, `email`, `password`, `phone`, `address`
**Output**: `accountId`
**Regras**:
- `email` deve ser únicos.
- Todos os campos obrigatórios devem ser válidos.

#### UpdateAccount (Atualizar Conta)
Atualizar informações da conta do psicólogo.
**Input**: `accountId`, `name`, `phone`, `address`, `defaultSessionValue`
**Output**: `void`
**Regras**:
- Apenas o titular da conta pode atualizar seus dados.
- `email` não pode ser alterado.

#### DeactivateAccount (Desativar Conta)
Desativar conta do psicólogo (soft delete).
**Input**: `accountId`, `reason`
**Output**: `void`
**Regras**:
- Altera o status da conta para `inactive`.
- Cancela todos os agendamentos futuros associados à conta.
- Os dados são mantidos para fins de auditoria.

---

### 🏥 Gestão de Convênios

#### CreateInsurance (Criar Convênio)
Cadastrar um novo convênio.
**Input**: `accountId`, `insuranceName`, `sessionValue`
**Output**: `insuranceId`
**Regras**: O nome do convênio deve ser único para a conta.

#### DeleteInsurance (Excluir Convênio)
Excluir um convênio (soft delete).
**Input**: `insuranceId`
**Output**: `void`
**Regras**:
- Não é possível excluir se houver pacientes ativos vinculados.
- Altera o status para `deleted`.

---

### Gestão de Pacientes

#### CreatePatient (Criar Paciente)
Cadastrar um novo paciente.
**Input**: `accountId`, `name`, `email`, `phone`, `birthDate`, `legalGuardian`
**Output**: `patientId`
**Regras**: `legalGuardian` é obrigatório se o paciente for menor de 18 anos.

#### UpdatePatient (Atualizar Paciente)
Atualizar informações de um paciente.
**Input**: `patientId`, `name`, `email`, `phone`, `status`
**Output**: `void`
**Regras**: Apenas o psicólogo responsável pode atualizar.

#### DeletePatient (Excluir Paciente)
Excluir um paciente (soft delete).
**Input**: `patientId`
**Output**: `void`
**Regras**:
- Altera o status do paciente para `deleted`.
- Cancela todos os seus agendamentos futuros.
- Histórico clínico é mantido para auditoria.

---

### Gestão de Agendamentos

#### ScheduleSession (Agendar Sessão)
Criar agendamento de sessão única ou recorrente.
**Input**: `accountId`, `patientId`, `dateTime`, `isRecurring`, `recurrenceType`
**Output**: `appointmentId` ou `appointmentIds[]`
**Regras**:
- Não pode haver conflito de horário na agenda da conta.
- Para recorrência, deve ser especificado um critério de término (data final ou número de ocorrências).

#### RescheduleSession (Remarcar Sessão)
Alterar data/hora de uma sessão.
**Input**: `appointmentId`, `newDateTime`
**Output**: `void`
**Regras**:
- A nova data/hora não pode ter conflito de horário.
- Se parte de uma série, apenas a instância atual é alterada.

#### CancelRecurringSeries (Cancelar Série Recorrente)
Cancelar todas as sessões futuras de uma série.
**Input**: `seriesId`
**Output**: `cancelledSessionsCount`
**Regras**: Cancela todas as sessões futuras (`scheduled`) da série.

---

### Gestão de Prontuários

#### CreateProgressNote (Criar Evolução)
Registrar evolução de uma sessão.
**Input**: `appointmentId`, `sessionSummary`
**Output**: `progressNoteId`
**Regras**: Apenas para sessões com status `completed`.

#### UpdateProgressNote (Atualizar Evolução)
Atualizar uma evolução.
**Input**: `progressNoteId`, `sessionSummary`
**Output**: `void`
**Regras**: A edição é bloqueada após 30 dias da criação, conforme normas do CFP.

---

### Relatórios

#### GetSchedule (Obter Agenda)
Retornar a agenda para um período.
**Input**: `accountId`, `startDate`, `endDate`
**Output**: `appointments[]`

#### GetFinancialReport (Obter Relatório Financeiro)
Retornar resumo financeiro para um período.
**Input**: `accountId`, `startDate`, `endDate`
**Output**: `totalRevenue`, `paidSessions`, `pendingSessions`

---

## Endpoints da API

```http
# Conta
POST   /signup
PUT    /profile
DELETE /profile

# Pacientes
POST   /patients
GET    /patients
PUT    /patients/:patientId
DELETE /patients/:patientId

# Convênios
POST   /insurances
GET    /insurances
DELETE /insurances/:insuranceId

# Agendamentos
POST   /appointments
PUT    /appointments/:appointmentId/reschedule
DELETE /appointments/:appointmentId
DELETE /recurring-series/:seriesId

# Prontuários
POST   /progress-notes
PUT    /progress-notes/:noteId

# Relatórios
GET    /schedule
GET    /reports/financial
```

---

## Modelo de Dados

### Tabela de Contas
```sql
CREATE TABLE psychological_management.account (
    account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    default_session_value DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
```

### Tabela de Convênios
```sql
CREATE TABLE psychological_management.insurance (
    insurance_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    name VARCHAR(255) NOT NULL,
    session_value DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (account_id, name)
);
```

### Tabela de Pacientes
```sql
CREATE TABLE psychological_management.patient (
    patient_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    birth_date DATE,
    legal_guardian_name VARCHAR(255),
    insurance_id UUID REFERENCES psychological_management.insurance(insurance_id) ON DELETE SET NULL,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'discharged', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
```

### Tabela de Séries Recorrentes
```sql
CREATE TABLE psychological_management.recurring_series (
    series_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    patient_id UUID NOT NULL,
    recurrence_type VARCHAR(20) NOT NULL CHECK (recurrence_type IN ('daily', 'weekly', 'monthly')),
    recurrence_interval INTEGER NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    max_occurrences INTEGER,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'cancelled', 'completed', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### Tabela de Agendamentos
```sql
CREATE TABLE psychological_management.appointment (
    appointment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES psychological_management.account(account_id),
    patient_id UUID NOT NULL REFERENCES psychological_management.patient(patient_id),
    series_id UUID REFERENCES psychological_management.recurring_series(series_id),
    scheduled_datetime TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER DEFAULT 50,
    session_value DECIMAL(10,2) NOT NULL,
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('private', 'insurance')),
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'no_show', 'cancelled')),
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### Tabela de Evoluções
```sql
CREATE TABLE psychological_management.progress_note (
    progress_note_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    appointment_id UUID NOT NULL REFERENCES psychological_management.appointment(appointment_id) UNIQUE,
    session_summary TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
```

---

## Architecture Decision Records

### ADR 001: Armazenar Dados Sensíveis com Criptografia
**Contexto**: O sistema lida com dados sensíveis de pacientes (prontuários, notas pessoais) que exigem conformidade com a LGPD e sigilo profissional.
**Decisão**: Implementar criptografia em nível de aplicação (ex: AES-256) para campos sensíveis, como o `session_summary` na tabela `progress_note`. As chaves de criptografia serão gerenciadas por um serviço seguro, separado do banco de dados.
**Consequências**:
- ✅ Maior segurança e conformidade legal.
- ❌ Aumenta a complexidade da aplicação e a latência de leitura/escrita desses dados.

### ADR 002: Geração de Sessões Recorrentes
**Contexto**: Sessões recorrentes podem gerar um grande volume de dados.
**Decisão**: Gerar todos os agendamentos de uma série recorrente no momento da sua criação e armazená-los individualmente na tabela `appointment`. A tabela `recurring_series` servirá como um registro mestre para agrupar essas sessões.
**Consequências**:
- ✅ Simplifica a consulta da agenda e permite modificações individuais em sessões de uma série.
- ❌ Aumenta o uso de armazenamento no banco de dados. Requer um limite máximo de ocorrências (ex: 52) para evitar sobrecarga.

### ADR 003: Soft Delete vs Hard Delete
**Contexto**: Dados de saúde exigem retenção para fins de auditoria, mesmo após a "exclusão" pelo usuário.
**Decisão**: Implementar soft delete usando um campo `status` e um timestamp `deleted_at` em todas as entidades principais. A lógica de negócio na aplicação será responsável por filtrar registros marcados como `deleted`.
**Consequências**:
- ✅ Mantém a rastreabilidade e permite a recuperação de dados.
- ❌ As queries da aplicação devem sempre incluir a cláusula `WHERE status != 'deleted'`, aumentando a complexidade e o risco de erro humano.

### ADR 004: Controle de Edição de Evoluções
**Contexto**: Prontuários clínicos têm regras rígidas de alteração, conforme o Conselho Federal de Psicologia (CFP).
**Decisão**: A lógica para controlar a edição de evoluções (`progress_note`) será implementada na camada de aplicação. Será verificado se a tentativa de edição ocorre dentro de 30 dias da criação do registro.
**Consequências**:
- ✅ Garante conformidade com as normas profissionais.
- ✅ Mantém a lógica de negócio na aplicação, facilitando testes e manutenções.

### ADR 005: Multi-Tenancy e Isolamento de Dados
**Contexto**: O sistema será usado por múltiplos psicólogos. Os dados de uma conta jamais devem ser acessíveis por outra.
**Decisão**: Implementar uma estratégia de Multi-Tenancy lógica. Todas as tabelas principais conterão uma coluna `account_id`. A camada de acesso a dados da aplicação será responsável por adicionar a cláusula `WHERE account_id = :current_account_id` a todas as consultas, garantindo o isolamento dos dados.
**Consequências**:
- ✅ Solução simples e de baixo custo, eficaz para o escopo do projeto.
- ❌ Depende criticamente da disciplina de implementação na aplicação para garantir a segurança. Requer testes rigorosos para prevenir vazamento de dados entre contas.
