import { exec, spawn } from 'child_process';
import { config } from './config';


const cmd = 'ts-node-dev'; 
const params = `--project .codedoc/tsconfig.json` 
            + ` -T --watch ${config.src.base},.codedoc`
            + ` --ignore-watch .codedoc/node_modules`
            + ` .codedoc/serve`;


if (process.platform === 'win32') {
  const child = exec(cmd + ' ' + params);

  child.stdout?.pipe(process.stdout);
  child.stderr?.pipe(process.stderr);
  child.on('close', () => {});
}
else {
  const child = spawn(cmd, [params], { stdio: 'inherit', shell: 'bash' });
  child.on('close', () => {});
}