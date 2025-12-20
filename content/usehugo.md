---
title: "Usehugo"
date: 2022-10-12T20:08:53+08:00
draft: false
---


# 如何使用hugo

```
import pygame  
import random  
import math  
import sys  
import string  
  
# 初始化 Pygamepygame.init()  
  
# 设置窗口尺寸  
WIDTH, HEIGHT = 1600, 900  
screen = pygame.display.set_mode((WIDTH, HEIGHT))  
pygame.display.set_caption("Python Fireworks with Letters")  
  
# 颜色定义  
BLACK = (0, 0, 0)  
WHITE = (255, 255, 255)  
COLORS = [  
    (255, 50, 50),  # 红色  
    (50, 255, 50),  # 绿色  
    (50, 50, 255),  # 蓝色  
    (255, 255, 50),  # 黄色  
    (255, 50, 255),  # 紫色  
    (50, 255, 255),  # 青色  
    (255, 150, 50),  # 橙色  
    (255, 255, 255)  # 白色  
]  
  
# 字母列表（可以自定义）  
LETTERS = ["H3C", "IMO"]  
# LETTERS =  ["H", "3", "C"]  
FONT_SIZE = 50  
HOLLOW_THICKNESS = 1  # 镂空字母的轮廓厚度  
  
DOT_SPACING = 2  # 1=连续，2=隔1点，3=隔2点，以此类推  
# ========== 新增：背景图加载 ==========def load_background_image(image_path, width, height):  
    """加载并缩放背景图"""  
    try:  
        # 加载图片  
        bg_image = pygame.image.load(image_path).convert_alpha()  
        # 缩放图片到窗口尺寸  
        bg_image = pygame.transform.scale(bg_image, (width, height))  
        return bg_image  
    except pygame.error as e:  
        print(f"无法加载背景图: {e}")  
        print("将使用黑色背景替代")  
        # 创建黑色背景  
        bg_image = pygame.Surface((width, height))  
        bg_image.fill(BLACK)  
        return bg_image  
  
# 加载背景图（请替换为你的背景图路径）  
# 可以使用：  
# 1. 绝对路径："C:/images/background.jpg"  
# 2. 相对路径："background.jpg"（图片放在代码同一目录）  
# 3. 留空使用黑色背景  
BACKGROUND_IMAGE_PATH = "bg.png"  # 修改这里为你的背景图路径  
background = load_background_image(BACKGROUND_IMAGE_PATH, WIDTH, HEIGHT)  
# 字母粒子类  
class LetterParticle:  
    def __init__(self, x, y, letter, color, velocity_x, velocity_y, gravity=0.1, decay=0.97):  
        self.x = x  
        self.y = y  
        self.letter = letter  
        self.color = color  
        self.velocity_x = velocity_x  
        self.velocity_y = velocity_y  
        self.gravity = gravity  
        self.decay = decay  
        self.lifetime = 255  # 用于淡出效果  
        self.rotation = random.uniform(0, 360)  
        self.rotation_speed = random.uniform(-2, 2)  
        self.rotate = random.random() < 0.8  
        if random.random() < 0.1:  
            self.font = pygame.font.SysFont(None, 150)  
            self.rotate = False  
        else:  
            self.font = pygame.font.SysFont(None, FONT_SIZE)  
        self.scale = 1.0  
        # 亮度变化  
        self.brightness = 255  
        self.brightness_speed = random.uniform(1, 2)  
        self.scale_decay = 0.99  
  
  
    def update(self):  
        self.x += self.velocity_x  
        self.y += self.velocity_y  
        self.velocity_y += self.gravity  
        self.velocity_x *= self.decay  
        self.velocity_y *= self.decay  
        self.lifetime -= 2  
        self.rotation += self.rotation_speed  
        self.scale *= self.scale_decay  
  
    def draw(self, surface):  
        if self.lifetime > 0 and self.brightness > 0:  
            # 计算当前透明度  
            alpha = max(0, min(255, self.lifetime * self.brightness / 255))  
  
            # 1. 绘制镂空字符  
            hollow_surface = self.draw_hollow_text(  
                surface, self.letter, (0, 0), self.color, self.font,  
                HOLLOW_THICKNESS, dot_spacing=2  # dot_spacing=2 是隔1个点绘制  
            )  
  
            # 2. 应用缩放  
            new_size = (int(hollow_surface.get_width() * self.scale),  
                        int(hollow_surface.get_height() * self.scale))  
            if new_size[0] > 0 and new_size[1] > 0:  
                hollow_surface = pygame.transform.scale(hollow_surface, new_size)  
  
            # 3. 应用旋转  
            if self.rotate:  
                hollow_surface = pygame.transform.rotate(hollow_surface, self.rotation)  
  
            # 4. 设置透明度  
            hollow_surface.set_alpha(alpha)  
  
            # 5. 绘制光晕效果（镂空字母的光晕更柔和）  
            # if alpha > 100:  
            #     glow_size = int(hollow_surface.get_width() * 1.1)            #            #     glow_surface = pygame.Surface((glow_size, glow_size), pygame.SRCALPHA)            #     # 绘制渐变光晕  
            #     for r in range(glow_size // 2, 0, -1):  
            #         glow_alpha = int(alpha * (r / (glow_size // 2)) / 20)  # 把 /10 改成 /20，让光晕更淡  
            #         # glow_alpha = int(alpha * (r / (glow_size // 2)) / 10)  
            #         if glow_alpha > 0:            #             pygame.draw.circle(            #                 glow_surface, (*self.color[:3], glow_alpha),            #                 (glow_size // 2, glow_size // 2), r            #             )            #     glow_rect = glow_surface.get_rect(center=(int(self.x), int(self.y)))            #     surface.blit(glow_surface, glow_rect)  
            # 6. 绘制镂空字符  
            text_rect = hollow_surface.get_rect(center=(int(self.x), int(self.y)))  
            surface.blit(hollow_surface, text_rect)  
  
    def is_dead(self):  
        return self.lifetime <= 0 or self.y > HEIGHT + 50 or self.scale < 0.3  
  
    def draw_hollow_text(self, surface, text, pos, color, font, thickness=2, dot_spacing=2):  
        """绘制点状镂空文字的核心方法  
        :param thickness: 轮廓厚度  
        :param dot_spacing: 点之间的间隔（越小越密，1=连续，2=隔1个点，3=隔2个点）  
        """        # 1. 先渲染实心文字作为掩码  
        text_surface = font.render(text, True, WHITE)  
        # 2. 获取文字的掩码（轮廓）  
        mask = pygame.mask.from_surface(text_surface)  
        # 3. 创建新的Surface用于绘制轮廓  
        hollow_surface = pygame.Surface(text_surface.get_size(), pygame.SRCALPHA)  
  
        # 4. 遍历掩码，绘制点状轮廓像素  
        for x in range(text_surface.get_width()):  
            for y in range(text_surface.get_height()):  
                if mask.get_at((x, y)):  
                    # 检查周围像素是否非掩码（轮廓边缘）  
                    is_edge = False  
                    # 上下左右检测  
                    for dx in [-1, 0, 1]:  
                        for dy in [-1, 0, 1]:  
                            if 0 <= x + dx < text_surface.get_width() and 0 <= y + dy < text_surface.get_height():  
                                if not mask.get_at((x + dx, y + dy)):  
                                    is_edge = True  
                                    break                        if is_edge:  
                            break  
  
                    # 只绘制轮廓边缘（点状效果）  
                    if is_edge:  
                        # 按间隔绘制点（x+y的和对间隔取模，形成均匀点状）  
                        if (x + y) % (dot_spacing + 1) == 0:  
                            # 绘制点状轮廓（可控制厚度）  
                            for dx in range(-thickness // 2, thickness // 2 + 1):  
                                for dy in range(-thickness // 2, thickness // 2 + 1):  
                                    if 0 <= x + dx < hollow_surface.get_width() and 0 <= y + dy < hollow_surface.get_height():  
                                        # 额外随机点效果（可选）  
                                        # if random.random() > 0.2:  # 随机跳过一些点，更自然  
                                        hollow_surface.set_at((x + dx, y + dy), color)  
  
        return hollow_surface  
  
  
# 烟花粒子类  
class Particle:  
    def __init__(self, x, y, color, velocity_x, velocity_y, radius=2, gravity=0.1, decay=0.97):  
        self.x = x  
        self.y = y  
        self.color = color  
        self.velocity_x = velocity_x  
        self.velocity_y = velocity_y  
        self.radius = radius  
        self.gravity = gravity  
        self.decay = decay  
        self.lifetime = 255  
  
    def update(self):  
        self.x += self.velocity_x  
        self.y += self.velocity_y  
        self.velocity_y += self.gravity  
        self.velocity_x *= self.decay  
        self.velocity_y *= self.decay  
        self.lifetime -= 2  
  
    def draw(self, surface):  
        if self.lifetime > 0:  
            alpha = max(0, self.lifetime)  
            pygame.draw.circle(surface, self.color,  
                               (int(self.x), int(self.y)),  
                               max(1, int(self.radius * self.lifetime / 255)))  
  
    def is_dead(self):  
        return self.lifetime <= 0 or self.y > HEIGHT + 10  
  
  
# 烟花类  
class Firework:  
    def __init__(self, x, y):  
        self.x = x  
        self.y = y  
        self.color = random.choice(COLORS)  
        self.velocity_y = random.uniform(-12, -8)  
        self.velocity_x = random.uniform(-1, 1)  
        self.gravity = 0.1  
        self.particles = []  
        self.letter_particles = []  # 新增：字母粒子列表  
        self.exploded = False  
        self.explosion_height = random.randint(100, 200)  
        self.explosion_power = random.randint(80, 150)  
        self.has_letters = random.random() < 0.8  # 50%的几率包含字母  
        # self.has_letters =True  
  
    def update(self):  
        if not self.exploded:  
            self.velocity_y += self.gravity  
            self.y += self.velocity_y  
            self.x += self.velocity_x  
  
            if random.random() < 0.3:  
                self.particles.append(  
                    Particle(self.x, self.y, self.color,  
                             random.uniform(-0.5, 0.5),  
                             random.uniform(-0.5, 0.5),  
                             radius=1.5, gravity=0.05, decay=0.9)  
                )  
  
            if self.velocity_y >= 0 or self.y <= self.explosion_height:  
                self.explode()  
        else:  
            for particle in self.particles[:]:  
                particle.update()  
                if particle.is_dead():  
                    self.particles.remove(particle)  
  
            # 更新字母粒子  
            for letter_particle in self.letter_particles[:]:  
                letter_particle.update()  
                if letter_particle.is_dead():  
                    self.letter_particles.remove(letter_particle)  
  
    def explode(self):  
        self.exploded = True  
        num_particles = self.explosion_power  
  
        # 创建普通粒子  
        for _ in range(num_particles):  
            angle = random.uniform(0, math.pi * 2)  
            speed = random.uniform(1, 5)  
            velocity_x = math.cos(angle) * speed  
            velocity_y = math.sin(angle) * speed  
  
            color_variation = random.randint(-30, 30)  
            color = (  
                min(255, max(0, self.color[0] + color_variation)),  
                min(255, max(0, self.color[1] + color_variation)),  
                min(255, max(0, self.color[2] + color_variation))  
            )  
  
            self.particles.append(  
                Particle(self.x, self.y, color, velocity_x, velocity_y,  
                         radius=random.uniform(1.5, 3.5),  
                         gravity=0.1, decay=0.96)  
            )  
  
        # 创建字母粒子（如果启用）  
        if self.has_letters:  
            # num_letters = random.randint(3, 8)  
            num_letters = 1  
            for _ in range(num_letters):  
                angle = random.uniform(0, math.pi * 2)  
                speed = random.uniform(0.5, 3)  
                velocity_x = math.cos(angle) * speed  
                velocity_y = math.sin(angle) * speed  
  
                letter = random.choice(LETTERS)  
  
                # 字母颜色可以稍微不同  
                letter_color = (  
                    min(255, max(0, self.color[0] + random.randint(-50, 50))),  
                    min(255, max(0, self.color[1] + random.randint(-50, 50))),  
                    min(255, max(0, self.color[2] + random.randint(-50, 50)))  
                )  
  
                self.letter_particles.append(  
                    LetterParticle(self.x, self.y, letter, letter_color,  
                                   velocity_x, velocity_y,  
                                   gravity=0.08, decay=0.95)  
                )  
  
    def draw(self, surface):  
        if not self.exploded:  
            pygame.draw.circle(surface, self.color, (int(self.x), int(self.y)), 3)  
  
        # 绘制普通粒子  
        for particle in self.particles:  
            particle.draw(surface)  
  
        # 绘制字母粒子  
        for letter_particle in self.letter_particles:  
            letter_particle.draw(surface)  
  
    def is_dead(self):  
        return self.exploded and len(self.particles) == 0 and len(self.letter_particles) == 0  
  
  
# 星星背景类（保持不变）  
class Star:  
    def __init__(self):  
        self.x = random.randint(0, WIDTH)  
        self.y = random.randint(0, HEIGHT)  
        self.size = random.uniform(0.1, 1.5)  
        self.brightness = random.randint(100, 255)  
        self.twinkle_speed = random.uniform(0.01, 0.05)  
        self.twinkle_offset = random.uniform(0, math.pi * 2)  
  
    def draw(self, surface, time):  
        twinkle = (math.sin(time * self.twinkle_speed + self.twinkle_offset) + 1) / 2  
        current_brightness = int(self.brightness * (0.7 + 0.3 * twinkle))  
        color = (current_brightness, current_brightness, current_brightness)  
        pygame.draw.circle(surface, color, (int(self.x), int(self.y)), self.size)  
  
  
# 主函数  
def main():  
    clock = pygame.time.Clock()  
    fireworks = []  
    stars = [Star() for _ in range(150)]  
    time = 0  
    font = pygame.font.SysFont(None, 36)  
  
    # 添加初始烟花  
    for _ in range(3):  
        fireworks.append(Firework(random.randint(100, WIDTH - 100), HEIGHT))  
  
    running = True  
    while running:  
        for event in pygame.event.get():  
            if event.type == pygame.QUIT:  
                running = False  
            elif event.type == pygame.KEYDOWN:  
                if event.key == pygame.K_ESCAPE:  
                    running = False  
                elif event.key == pygame.K_SPACE:  
                    for _ in range(5):  
                        fireworks.append(Firework(random.randint(100, WIDTH - 100), HEIGHT))  
                elif event.key == pygame.K_l:  # L键：强制生成带字母的烟花  
                    for _ in range(3):  
                        fw = Firework(random.randint(100, WIDTH - 100), HEIGHT)  
                        fw.has_letters = True  # 强制包含字母  
                        fireworks.append(fw)  
            elif event.type == pygame.MOUSEBUTTONDOWN:  
                if event.button == 1:  
                    fireworks.append(Firework(event.pos[0], HEIGHT))  
                elif event.button == 3:  # 右键：生成带字母的烟花  
                    fw = Firework(event.pos[0], HEIGHT)  
                    fw.has_letters = True  
                    fireworks.append(fw)  
  
        # 随机添加新烟花  
        if random.random() < 0.03 and len(fireworks) < 20:  
            fireworks.append(Firework(random.randint(100, WIDTH - 100), HEIGHT))  
  
        # 更新烟花  
        for firework in fireworks[:]:  
            firework.update()  
            if firework.is_dead():  
                fireworks.remove(firework)  
  
        # 更新星星  
        time += 0.05  
        # ========== 修改：绘制背景图 ==========        # 先绘制背景图  
        screen.blit(background, (0, 0))  
        # 如果需要背景图半透明效果，可以添加：  
        # overlay = pygame.Surface((WIDTH, HEIGHT), pygame.SRCALPHA)  
        # overlay.fill((0, 0, 0, 100))  # 黑色半透明遮罩，最后一个值是透明度（0-255）  
        # screen.blit(overlay, (0, 0))  
        # 绘制背景  
        # screen.fill(BLACK)  
  
        # 绘制星星  
        for star in stars:  
            star.draw(screen, time)  
  
        # 绘制烟花  
        for firework in fireworks:  
            firework.draw(screen)  
  
        # 显示说明文字  
        instructions = [  
            "Left Click: Add Firework",  
            "Right Click: Add Firework with Letters",  
            "Space: Add 5 Fireworks",  
            "L: Add 3 Letter Fireworks",  
            "ESC: Exit"  
        ]  
  
        # for i, text in enumerate(instructions):  
        #     text_surface = font.render(text, True, WHITE)        #     screen.blit(text_surface, (10, 10 + i * 30))  
        # 显示烟花和字母数量  
        letter_count = sum(len(fw.letter_particles) for fw in fireworks)  
        count_text = f"Fireworks: {len(fireworks)}  Letters: {letter_count}"  
        count_surface = font.render(count_text, True, WHITE)  
        # screen.blit(count_surface, (WIDTH - 300, 10))  
  
        pygame.display.flip()  
        clock.tick(60)  
  
    pygame.quit()  
    sys.exit()  
  
  
if __name__ == "__main__":  
    main()
```