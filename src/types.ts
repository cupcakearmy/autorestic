export type StringOrArray = string | string[]

// BACKENDS

type BackendLocal = {
	type: 'local'
	key: string
	path: string
}

type BackendSFTP = {
	type: 'sftp'
	key: string
	path: string
	password?: string
}

type BackendREST = {
	type: 'rest'
	key: string
	path: string
	user?: string
	password?: string
}

type BackendS3 = {
	type: 's3'
	key: string
	path: string
	aws_access_key_id: string
	aws_secret_access_key: string
}

type BackendB2 = {
	type: 'b2'
	key: string
	path: string
	b2_account_id: string
	b2_account_key: string
}

type BackendAzure = {
	type: 'azure'
	key: string
	path: string
	azure_account_name: string
	azure_account_key: string
}

type BackendGS = {
	type: 'gs'
	key: string
	path: string
	google_project_id: string
	google_application_credentials: string
}

export type Backend =
	| BackendAzure
	| BackendB2
	| BackendGS
	| BackendLocal
	| BackendREST
	| BackendS3
	| BackendSFTP

export type Backends = { [name: string]: Backend }

// LOCATIONS

export type Location = {
	from: string
	to: StringOrArray
	cron?: string
	hooks?: {
		before?: StringOrArray
		after?: StringOrArray
	}
	options?: {
		[key: string]: {
			[key: string]: StringOrArray
		}
	}
}

export type Locations = { [name: string]: Location }

// OTHER

export type Config = {
	locations: Locations
	backends: Backends
}

export type Lockfile = {
	running: boolean
	crons: {
		[name: string]: {
			lastRun: number
		}
	}
}

export type Flags = { [arg: string]: any }
