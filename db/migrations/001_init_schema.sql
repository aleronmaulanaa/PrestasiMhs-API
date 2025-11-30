-- 1. Aktifkan Extension UUID (Wajib untuk tipe data UUID)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Tabel Roles (SRS 3.1.2)
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert Data Awal Roles (SRS 3.1.2 Data Awal)
INSERT INTO roles (name, description) VALUES 
('Admin', 'Pengelola sistem'),
('Mahasiswa', 'Pelapor prestasi'),
('Dosen Wali', 'Verifikator prestasi');

-- 3. Tabel Users (SRS 3.1.1)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID REFERENCES roles(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 4. Tabel Permissions (SRS 3.1.3)
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

-- 5. Tabel Role Permissions (SRS 3.1.4)
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id),
    permission_id UUID REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

-- 6. Tabel Lecturers / Dosen (SRS 3.1.6)
CREATE TABLE lecturers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    lecturer_id VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 7. Tabel Students / Mahasiswa (SRS 3.1.5)
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    student_id VARCHAR(20) UNIQUE NOT NULL,
    program_study VARCHAR(100),
    academic_year VARCHAR(10),
    advisor_id UUID REFERENCES lecturers(id), -- Relasi ke Dosen Wali
    created_at TIMESTAMP DEFAULT NOW()
);

-- 8. Enum Type untuk Status Prestasi (SRS 3.1.7)
CREATE TYPE achievement_status AS ENUM ('draft', 'submitted', 'verified', 'rejected', 'deleted');
-- Catatan: 'deleted' ditambahkan sesuai Highlight Aturan No. 2 (Soft Delete)

-- 9. Tabel Achievement References (SRS 3.1.7)
CREATE TABLE achievement_references (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    mongo_achievement_id VARCHAR(24) NOT NULL, -- Relasi ke MongoDB ID
    status achievement_status DEFAULT 'draft',
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id),
    rejection_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);