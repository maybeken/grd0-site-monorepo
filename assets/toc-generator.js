import fs from 'fs';
import path from 'path';

const output_filename = 'output.json';
const skip_file = [import.meta.filename.split('/').reverse()[0], output_filename];

async function getAllFilesRecursively(directory) {
  let files = [];

  try {
    const childFiles = await fs.promises.readdir(directory);

    for (let file of childFiles) {
      const fullPath = path.join(directory, file);
      const relativePath = fullPath.replace(import.meta.dirname, '');

      if (fs.statSync(fullPath).isDirectory()) {
        files = files.concat(await getAllFilesRecursively(fullPath));
      } else if (!(skip_file.includes(file))) {
        files.push(relativePath)
      }
    }
  } catch (err) {
    console.error("Error: ", err);
  }

  return files;
}

async function groupFilesByDirectory(files) {
  const groups = files.map((val) => {
    const path_split = val.split('/');
    const path = path_split.slice(1).reverse().slice(1);
    const filename = path_split.reverse()[0];

    return {
      path: path ? `/${path.reverse().join('/')}` : '/',
      filename,
    }
  });

  const group_by = Object.groupBy(groups, ({ path }) => path);

  return flattenGroups(group_by);
}

async function flattenGroups(groups) {
  for (let key in groups) {
    groups[key] = groups[key].map((val) => {
      return val.filename;
    });
  }

  return groups;
}

async function main() {
  let currentDir = process.cwd();
  let allFiles = await getAllFilesRecursively(currentDir);

  let grouped = await groupFilesByDirectory(allFiles);

  fs.writeFile(output_filename, JSON.stringify(grouped), 'utf8', (err) => {
    if (err) {
      console.log("Error writing file");
      return;
    }
  });
}

main().catch(console.error);
