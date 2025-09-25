import { Routes } from '@angular/router';
import { Search } from '../search/search';
import { DriveManagement } from './drive-management.routes';
import { FilesComponent } from './files/files/files';
import { FoldersComponent } from './folders/folders/folders';
import { StatsComponent } from './stats/stats/stats';
import { Upload } from './upload/upload/upload';
export const driveManagementRoutes: Routes = [
  {
    path: '',
    component: DriveManagement,
    children: [
      { path: 'folders', component: FoldersComponent },
      { path: 'files', component: FilesComponent },
      { path: 'upload', component: Upload },
      { path: 'search', component: Search },
      { path: 'stats', component: StatsComponent },
      { path: '', redirectTo: 'folders', pathMatch: 'full' },
    ],
  },
];
