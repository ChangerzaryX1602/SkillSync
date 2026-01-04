import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { UserService } from '../../core/services/user.service';

@Component({
  selector: 'app-home',
  imports: [],
  templateUrl: './home.html',
  styleUrl: './home.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HomeComponent implements OnInit {
  private readonly userService = inject(UserService);

  protected readonly skills = [
    { name: 'TypeScript', level: 85, color: '#3178c6' },
    { name: 'Angular', level: 78, color: '#dd0031' },
    { name: 'Golang', level: 72, color: '#00add8' },
    { name: 'PostgreSQL', level: 68, color: '#336791' },
    { name: 'Docker', level: 65, color: '#2496ed' },
  ];

  protected readonly courses = [
    {
      id: 1,
      title: 'Angular ขั้นสูง (Advanced Patterns)',
      instructor: 'สมชาย โค้ดเดอร์',
      progress: 65,
      thumbnail: 'https://picsum.photos/seed/angular/300/200',
      category: 'Frontend',
    },
    {
      id: 2,
      title: 'สร้าง Microservices ด้วย Go',
      instructor: 'สมหญิง แบ็คเอนด์',
      progress: 40,
      thumbnail: 'https://picsum.photos/seed/golang/300/200',
      category: 'Backend',
    },
    {
      id: 3,
      title: 'พื้นฐาน System Design (ฉบับวัยรุ่น)',
      instructor: 'อเล็กซ์ สถาปนิก',
      progress: 20,
      thumbnail: 'https://picsum.photos/seed/system/300/200',
      category: 'Architecture',
    },
  ];

  protected readonly stats = [
    { label: 'คอร์สที่ลงเรียน', value: 12, icon: '📚' },
    { label: 'ทักษะที่เชี่ยวชาญ', value: 8, icon: '🎯' },
    { label: 'ชั่วโมงเรียน', value: 156, icon: '⏱️' },
    { label: 'เกียรติบัตร', value: 5, icon: '🏆' },
  ];

  ngOnInit(): void {
    this.userService.fetchCurrentUser().subscribe();
  }
}
